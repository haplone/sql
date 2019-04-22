package main

import (
	"encoding/json"
	"fmt"
	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/meta"
	"github.com/pingcap/tidb/session"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/tablecodec"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {

	getRegions()
	getTbls()
}

func getTbls() {
	err := session.RegisterStore("tikv", tikv.Driver{})
	check(err)
	store, err := session.NewStore("tikv://127.0.0.1:2379?disableGC=true")
	check(err)
	//se, err := session.CreateSession(store)
	//check(err)
	//se.NewTxn()

	kv.RunInNewTxn(store, false, func(txn kv.Transaction) error {
		t := meta.NewMeta(txn)
		dbs, err := t.ListDatabases()
		check(err)
		for _, db := range dbs {
			log.Printf("db : %s[%d]\n", db.Name, db.ID)

			tbls, err := t.ListTables(db.ID)
			check(err)
			for _, tbl := range tbls {
				prefix := tablecodec.EncodeTablePrefix(tbl.ID)
				prefix2 := strings.Trim(fmt.Sprintf("%q", prefix), "\"")
				//tablecodec.EncodeTablePrefix(prefix)
				log.Printf("db : %s.%s[%d.%d] key: %s\n", db.Name, tbl.Name, db.ID, tbl.ID, prefix2)
				//tbl.Indices
			}
		}
		return nil
	})

	store.Close()
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
func getRegions() {
	req, err := getRequest("", "", "", nil)
	if err != nil {
		log.Println(err)
	}
	dialClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: nil,
	}}
	res, err := dialClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	r, err := ioutil.ReadAll(res.Body)
	//log.Println(string(r))
	regions := regionsInfo{}
	err = json.Unmarshal(r, &regions)
	if err != nil {
		log.Println(err)
	}

	for _, rr := range regions.Regions {
		log.Println(rr)
	}

}
func getRequest(prefix string, method string, bodyType string, body io.Reader) (*http.Request, error) {
	if method == "" {
		method = http.MethodGet
	}
	url := "http://127.0.0.1:2379/pd/api/v1/regions"
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return req, err
}

type regionInfo struct {
	ID          uint64              `json:"id"`
	StartKey    string              `json:"start_key"`
	EndKey      string              `json:"end_key"`
	RegionEpoch *metapb.RegionEpoch `json:"epoch,omitempty"`
	Peers       []*metapb.Peer      `json:"peers,omitempty"`

	Leader          *metapb.Peer      `json:"leader,omitempty"`
	DownPeers       []*pdpb.PeerStats `json:"down_peers,omitempty"`
	PendingPeers    []*metapb.Peer    `json:"pending_peers,omitempty"`
	WrittenBytes    uint64            `json:"written_bytes,omitempty"`
	ReadBytes       uint64            `json:"read_bytes,omitempty"`
	ApproximateSize int64             `json:"approximate_size,omitempty"`
}

type regionsInfo struct {
	Count   int           `json:"count"`
	Regions []*regionInfo `json:"regions"`
}
