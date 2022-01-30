package plugin_callback

import (
	"github.com/Monibuca/engine/v3"
	. "github.com/Monibuca/utils/v3"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"io/ioutil"
	"log"
	"net/http"
)

type CallbackSRS struct {
	Action   string `json:"action"`
	App      string `json:"app"`
	ClientID uint32 `json:"client_id"`
	IP       string `json:"ip"`
	Param    string `json:"param"`
	Stream   string `json:"stream"`
	Vhost    string `json:"vhost"`
}

var config = struct {
	Debug     bool
	Connect   string
	Publish   string
	UnPublish string
	Close     string
	Play      string
	Stop      string
}{}

func init() {
	pc := engine.PluginConfig{
		Name:   "CALLBACK",
		Config: &config,
	}
	pc.Install(run)

	http.HandleFunc("/callback/list", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		data, _ := gjson.Encode(config)
		w.Write(data)
	})
	http.HandleFunc("/callback/test", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)

		body := []byte("")
		body, _ = ioutil.ReadAll(r.Body)
		w.Write(body)
	})
}

func run() {
	go engine.AddHookGo(engine.HOOK_PUBLISH, callbackPublish)
	go engine.AddHookGo(engine.HOOK_STREAMCLOSE, callbackCloses)
}

func callbackPublish(s *engine.Stream) {
	call := CallbackSRS{
		Action: "publish",
	}
	call.callback(s, config.Publish)
}

func callbackCloses(s *engine.Stream) {
	call := CallbackSRS{
		Action: "unpublish",
	}
	call.callback(s, config.Publish)
}

func (c *CallbackSRS) callback(s *engine.Stream, url string) {
	if g.IsEmpty(url) {
		glog.Infof("callback %s's url is empty", c.Action)
		return
	}

	param := gstr.Explode("/", s.StreamPath)
	if len(param) < 2 {
		return
	}

	c.App = param[0]
	c.Stream = param[len(param)-1]

	go func() {
		response, err := g.Client().ContentJson().Post(url, c)
		if err != nil {
			glog.Warning(err)
		}
		if config.Debug {
			response.RawDump()
		}

		log.Printf("CALLBACK %s/%s -> %s:%d", c.App, c.Stream, c.Action, response.StatusCode)
	}()
}
