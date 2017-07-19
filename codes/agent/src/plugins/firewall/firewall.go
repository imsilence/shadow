package firewall

const name string = "firewall"

type FirewallPlugin struct {
    cxt *config.Context
}

func (h *FirewallPlugin) Init(cxt *config.Context) (err error) {
    h.cxt = cxt
    h.cxt.Channels.Log <- fmt.Sprintf("Init Plugin %s", name)
    return nil
}

func (h *FirewallPlugin) Run() (err error) {
    return nil
}


func (h *FirewallPlugin) Call(command *entity.Command) (result *entity.CommandResult, err error) {
    return nil, nil
}

func (h *FirewallPlugin) Destory() (err error) {
    return nil
}

func init() {
    plugins.Register(name, new(FirewallPlugin))
}
