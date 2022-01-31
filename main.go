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
	"time"
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
	Timeout   int
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
		Name:   "Callback",
		Config: &config,
		HotConfig: map[string]func(interface{}){
			"Timeout": func(v interface{}) {
				config.Timeout = v.(int)
			},
			"Debug": func(v interface{}) {
				config.Debug = v.(bool)
			},
		},
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
	call.callback(s.StreamPath, config.Publish)
}

// engine 为了减少内存泄露，改成了 StreamPath
// https://github.com/Monibuca/engine/blob/c9052a981a31516f45690f1a3d36c36441fab4e2/stream.go#L131
func callbackCloses(StreamPath string) {
	call := CallbackSRS{
		Action: "unpublish",
	}
	call.callback(StreamPath, config.Publish)
}

func (c *CallbackSRS) callback(StreamPath, url string) {
	if g.IsEmpty(url) {
		glog.Infof("callback %s's url is empty", c.Action)
		return
	}

	param := gstr.Explode("/", StreamPath)
	if len(param) < 2 {
		return
	}

	c.App = param[0]
	c.Stream = param[len(param)-1]

	go func() {
		if config.Timeout < 0 || config.Timeout > 100 {
			config.Timeout = 30
		}

		response, err := g.Client().Timeout(time.Second*time.Duration(config.Timeout)).ContentJson().Post(url, c)
		if err != nil {
			glog.Warning(err)
		}
		if config.Debug {
			response.RawDump()
		}

		log.Printf("Plugin Callback %s/%s -> %s:%d", c.App, c.Stream, c.Action, response.StatusCode)
	}()
}
