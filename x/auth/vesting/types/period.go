package types

import (
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Periods stores all vesting periods passed as part of a PeriodicVestingAccount
type Periods []Period

// String Period implements stringer interface
func (p Period) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// TotalLength return the total length in seconds for a period
func (p Periods) TotalLength() int64 {
	var total int64
	for _, period := range p {
		total += period.Length
	}
	return total
}

// TotalDuration returns the total duration of the period
func (p Periods) TotalDuration() time.Duration {
	len := p.TotalLength()
	return time.Duration(len) * time.Second
}

// TotalDuration returns the sum of coins for the period
func (p Periods) TotalAmount() sdk.Coins {
	total := sdk.Coins{}
	for _, period := range p {
		total = total.Add(period.Amount...)
	}
	return total
}

// String Periods implements stringer interface
func (vp Periods) String() string {
	periodsListString := make([]string, len(vp))
	for _, period := range vp {
		periodsListString = append(periodsListString, period.String())
	}

	return strings.TrimSpace(fmt.Sprintf(`Vesting Periods:
		%s`, strings.Join(periodsListString, ", ")))
}
