package goonep

import (
   "encoding/json"
   "log"
)

type Pdevice struct{

   Basic    struct{   
      Modified  float32 `json:"modified,omitempty"`
      Status    string `json:"status,omitempty"`
      Subscribers   int `json:"subscribers,omitempty"`
      Type  string `json:"type,omitempty"`
   } `json:"basic,omitempty"`

   Comments []interface{} `json:"comments,omitempty"`

   Counts   struct{   
      Client     int `json:"client,omitempty"`
      Dataport   int `json:"dataport,omitempty"`
      Datarule   int `json:"datarule,omitempty"`
      Disk      int `json:"disk,omitempty"`
      Dispatch   int `json:"dispatch,omitempty"`
      Email     int `json:"email,omitempty"`
      Http      int `json:"http,omitempty"`
      Share     int `json:"share,omitempty"`
      Sms       int `json:"sms,omitempty"`
      Xmpp      int `json:"xmpp,omitempty"`
   } `json:"counts,omitempty"`

   Description  struct {   
      Limits    struct {   
         Client     int `json:"client,omitempty"`
         Dataport       string `json:"dataport,omitempty"`
         Datarule       string `json:"datarule,omitempty"`
         Disk          string `json:"disk,omitempty"`
         Dispatch       string `json:"dispatch,omitempty"`
         Email         string `json:"email,omitempty"`
         Email_bucket string `json:"email_bucket,omitempty"`
         Http         string `json:"http,omitempty"`
         Http_bucket     string `json:"http_bucket,omitempty"`
         Share         string `json:"share,omitempty"`
         Sms           string `json:"sms,omitempty"`
         Sms_bucket  string `json:"sms_bucket,omitempty"`
         Xmpp          string `json:"xmpp,omitempty"`
         Xmpp_bucket     string `json:"xmpp_bucket,omitempty"`
      } `json:"limits,omitempty"`
      Locked    bool `json:"locked,omitempty"`
      Meta     string `json:"meta,omitempty"`
      Name     string `json:"name,omitempty"`
      Public    bool `json:"public,omitempty"`
   } `json:"description,omitempty"`

   Subscribers  []interface{} `json:"subscribers,omitempty"`

   Tags []interface{} `json:"tags,omitempty"`

   Usage    struct {   
      Client    int `json:"client,omitempty"`
      Dataport  int `json:"dataport,omitempty"`
      Datarule  int `json:"datarule,omitempty"`
      Disk      int `json:"disk,omitempty"`
      Dispatch  int `json:"dispatch,omitempty"`
      Email     int `json:"email,omitempty"`
      Http      int `json:"http,omitempty"`
      Share     int `json:"share,omitempty"`
      Sms       int `json:"sms,omitempty"`
      Xmpp      int `json:"xmpp,omitempty"`
   } `json:"usage,omitempty"`

}

// Validate retrived device is valid or not
func (d Pdevice) Validate() bool{
   if len( d.Description.Meta ) <= 0 { return false }

   return true
}

// Meta from the meta field in devices' description information
func (d *Pdevice) GetMeta() DeviceMeta {

   meta := DeviceMeta{}

   err := json.Unmarshal( []byte( d.Description.Meta ) , &meta )

   if err != nil {
      log.Printf("Unmarshal Device [%s] meta met wrong: %v", d.Description.Name, err)
   }

   return meta

}

// Set meta for device
func (d *Pdevice) SetMeta(meta DeviceMeta) *Pdevice {

   metaString, err := json.Marshal( meta )

   if err != nil {
      log.Printf("Marshal meta for Device [%s] meta met wrong: %v", d.Description.Name, err)
   }

   d.Description.Meta = string(metaString)

   return d
}

type DeviceMeta struct {

   DeviceType      string `json:"deviceType",omitempty`
   DeviceTypeID     string `json:"deviceTypeID",omitempty`
   DeviceTypeName   string `json:"deviceTypeName",omitempty`
   Location       string `json:"Location",omitempty`
   Timezone       string `json:"Timezone",omitempty`
   Activetime      string `json:"activetime",omitempty`

   Device struct { 
      Model    string `json:"model",omitempty`
      Sn          string `json:"sn",omitempty`
      Type     string `json:"type",omitempty`
      Vendor    string `json:"vendor",omitempty`
   } `json:"device",omitempty`
   
   ExtraField string `json:"extra_field",omitempty`

}