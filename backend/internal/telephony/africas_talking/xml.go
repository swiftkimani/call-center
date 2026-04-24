package africas_talking

import (
	"fmt"
)

// BuildInboundResponse returns Africa's Talking action XML that bridges the call.
// agentPhone is the device/SIP URI to ring; callbackURL receives status events.
func (a *Adapter) BuildInboundResponse(agentPhone, callbackURL string) ([]byte, error) {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Dial record="true" sequential="false" callbackUrl="%s">
    <Number>%s</Number>
  </Dial>
</Response>`, callbackURL, agentPhone)
	return []byte(xml), nil
}
