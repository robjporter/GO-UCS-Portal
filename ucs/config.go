package ucs

import (
    "io"
    "os"
    "bytes"
    "errors"
    "encoding/json"
)

func LoadConfig(configFile string) error {
    if _, err := os.Stat(configFile); err == nil {
        buf := bytes.NewBuffer(nil)
        f, _ := os.Open(configFile) // Error handling elided for brevity.
        io.Copy(buf, f)           // Error handling elided for brevity.
        f.Close()
        var config map[string]interface{}
        if err := json.Unmarshal(buf.Bytes(), &config); err != nil {
            return err
        }
        _, ok1 := config["ip"]
        _, ok2 := config["username"]
        _, ok3 := config["password"]
        if ok1 && ok2 && ok3 {
            SetUCSAddress(config["ip"].(string))
            SetUCSUsername(config["username"].(string))
            SetUCSPassword(config["password"].(string))
            return nil
        } else {
            return errors.New("The config file could not be read or is invalid.  Please visit /config to update.")
        }
    }
    return errors.New("The config file could not be found.  Please ensure it exists inside the config directory.")

}
