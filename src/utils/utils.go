package utils

import (
    "crypto/rand"
    "encoding/json"
    "errors"
    "fmt"
    "net"
    "net/url"
    "reflect"
    "strconv"
    "strings"
    "time"
    "log"
)

// FormatEndpoint formats endpoint
func FormatEndpoint(endpoint string) string {
    endpoint = strings.TrimSpace(endpoint)
    endpoint = strings.TrimRight(endpoint, "/")
    if !strings.HasPrefix(endpoint, "http://") &&
        !strings.HasPrefix(endpoint, "https://") {
        endpoint = "http://" + endpoint
    }

    return endpoint
}

// ParseEndpoint parses endpoint to a URL
func ParseEndpoint(endpoint string) (*url.URL, error) {
    endpoint = FormatEndpoint(endpoint)

    u, err := url.Parse(endpoint)
    if err != nil {
        return nil, err
    }
    return u, nil
}

// ParseRepository splits a repository into two parts: project and rest
func ParseRepository(repository string) (project, rest string) {
    repository = strings.TrimLeft(repository, "/")
    repository = strings.TrimRight(repository, "/")
    if !strings.ContainsRune(repository, '/') {
        rest = repository
        return
    }
    index := strings.Index(repository, "/")
    project = repository[0:index]
    rest = repository[index + 1:]
    return
}

// GenerateRandomString generates a random string
func GenerateRandomString() string {
    length := 32
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
    l := len(chars)
    result := make([]byte, length)
    _, err := rand.Read(result)
    if err != nil {
        fmt.Errorf("Error reading random bytes: %v", err)
    }
    for i := 0; i < length; i++ {
        result[i] = chars[int(result[i]) % l]
    }
    return string(result)
}

// TestTCPConn tests TCP connection
// timeout: the total time before returning if something is wrong
// with the connection, in second
// interval: the interval time for retring after failure, in second
func TestTCPConn(addr string, timeout, interval int) error {
    success := make(chan int)
    cancel := make(chan int)

    go func() {
        for {
            select {
            case <-cancel:
                break
            default:
                conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout) * time.Second)
                if err != nil {
                    log.Printf("failed to connect to tcp://%s, retry after %d seconds :%v",
                        addr, interval, err)
                    time.Sleep(time.Duration(interval) * time.Second)
                    continue
                }
                conn.Close()
                success <- 1
                break
            }
        }
    }()

    select {
    case <-success:
        return nil
    case <-time.After(time.Duration(timeout) * time.Second):
        cancel <- 1
        return fmt.Errorf("failed to connect to tcp:%s after %d seconds", addr, timeout)
    }
}

// ParseTimeStamp parse timestamp to time
func ParseTimeStamp(timestamp string) (*time.Time, error) {
    i, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return nil, err
    }
    t := time.Unix(i, 0)
    return &t, nil
}

//ConvertMapToStruct is used to fill the specified struct with map.
func ConvertMapToStruct(object interface{}, values interface{}) error {
    if object == nil {
        return errors.New("nil struct is not supported")
    }

    if reflect.TypeOf(object).Kind() != reflect.Ptr {
        return errors.New("object should be referred by pointer")
    }

    bytes, err := json.Marshal(values)
    if err != nil {
        return err
    }

    return json.Unmarshal(bytes, object)
}

// ParseProjectIDOrName parses value to ID(int64) or name(string)
func ParseProjectIDOrName(value interface{}) (int64, string, error) {
    if value == nil {
        return 0, "", errors.New("harborIDOrName is nil")
    }

    var id int64
    var name string
    switch value.(type) {
    case int:
        i := value.(int)
        id = int64(i)
        if id == 0 {
            return 0, "", fmt.Errorf("invalid ID: 0")
        }
    case int64:
        id = value.(int64)
        if id == 0 {
            return 0, "", fmt.Errorf("invalid ID: 0")
        }
    case string:
        name = value.(string)
        if len(name) == 0 {
            return 0, "", fmt.Errorf("empty name")
        }
    default:
        return 0, "", fmt.Errorf("unsupported type")
    }
    return id, name, nil
}

