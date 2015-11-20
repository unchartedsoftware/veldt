package elastic

/*
import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"encoding/binary"
	"encoding/json"

	"github.com/parnurzeal/gorequest"

	"github.com/unchartedsoftware/prism/binning"
)
*/

const esHost = "http://10.64.16.120:9200"
const esIndex = "nyc_twitter"
const tileResolution = 256
const maxLevelSupported = 24
