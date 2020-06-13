# Jitsi MUC Exporter

Prometheus exporter for Jitsi that works by listening on JvbBrewery MUC room.

> WARNING! I've only added support for the stats that I care about in my installation.  
  Contributions are welcome if you want to add support for more stats and fill the help texts.

## Why MUC instead of REST API?

**TL;DR**: Deploy one exporter per cluster instead of many (one per JVB)

Jitsi Videobridge can expose its stats with a REST API. This is fine, and you can use something like [Telegraf](https://github.com/influxdata/telegraf)
to transform the JSON response from the API to a set of usable Prometheus metrics.

However, if you start scaling your Jitsi installation by adding more JVBs, you will have to deploy your exporters
for each installation. With this exporter you don't have to.

All JVBs also report their stats via XMPP to a JvbBrewery MUC room.
Jicofo listens for those stats in order to distribute meetings between instances of JVBs.  
This means there's a central place where you can get stats about **all** your JVB instances,
so you only have to deploy the Prometheus exporter once.

This makes it easier to deal with service discovery (or more precisely - a lack of it) and static targets configurations.

## Contributions

All contributions are welcome, especially those that add **support for missing stats**.

## TODO

- [ ] Add support for all stats
- [ ] Improve help texts for stats
- [ ] Configuration with env vars
- [ ] Create a Telegraf plugin for Jitsi based on the principals of this project
- [ ] Create a Grafana dashboard template
