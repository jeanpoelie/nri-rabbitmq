package args

import (
	"regexp"
	"testing"

	"github.com/newrelic/nri-rabbitmq/utils/consts"
	"github.com/stretchr/testify/assert"
)

func TestParseFilterArgsDefault(t *testing.T) {
	var argList = ArgumentList{}
	err := SetGlobalArgs(argList)
	assert.NoError(t, err, "err should be nil")
}

func TestParseFilterArgsBadArg(t *testing.T) {
	var argList = ArgumentList{
		Exchanges: "invalid",
	}
	err := SetGlobalArgs(argList)
	assert.Error(t, err, "should have error from bad arg when unmarshaling")

	argList.Exchanges = ""
	argList.ExchangesRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.ExchangesRegexes = ""
	argList.Queues = "invalid"
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.Queues = ""
	argList.QueuesRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.QueuesRegexes = ""
	argList.Vhosts = "invalid"
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.Vhosts = ""
	argList.VhostsRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.VhostsRegexes = `["(invalid-group"]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)
}

func TestParseFilterArgsValidJson(t *testing.T) {
	var argList = ArgumentList{
		Queues:        `["test-1", "test-2", "test-3"]`,
		QueuesRegexes: `["one-.*", "two-.*"]`,
	}
	err := SetGlobalArgs(argList)
	assert.NoError(t, err)
	assert.NotNil(t, GlobalArgs)
	assert.Equal(t, 3, len(GlobalArgs.Queues))
	assert.Equal(t, 2, len(GlobalArgs.QueuesRegexes))
	assert.True(t, GlobalArgs.QueuesRegexes[0].MatchString("one-queue"))
	assert.False(t, GlobalArgs.QueuesRegexes[0].MatchString("two-queue"))
	assert.False(t, GlobalArgs.QueuesRegexes[1].MatchString("one-queue"))
	assert.True(t, GlobalArgs.QueuesRegexes[1].MatchString("two-queue"))
}

func TestCheckArgsBadInventory(t *testing.T) {
	testArgs := RabbitMQArguments{}
	testArgs.Metrics = true
	testArgs.Inventory = false
	err := testArgs.Validate()
	assert.Error(t, err, "To collect Metrics, you also need to collect Inventory")
}

func TestCheckArgsBadMetrics(t *testing.T) {
	testArgs := RabbitMQArguments{}
	testArgs.Metrics = false
	testArgs.Inventory = false
	testArgs.Events = true
	err := testArgs.Validate()
	assert.Error(t, err, "It should collect at lest Inventory, or Inventory and Metrics")
}

func TestCheckArgsGoodArgs(t *testing.T) {
	testArgs := RabbitMQArguments{}
	err := testArgs.Validate()
	assert.NoError(t, err, "Collecting everything (no args) should be valid")

	testArgs.Inventory = true
	err = testArgs.Validate()
	assert.NoError(t, err, "Collecting just Inventory (-inventory) should be valid")

	testArgs.Metrics = true
	err = testArgs.Validate()
	assert.NoError(t, err, "Collecting both Inventory and Metrics (-inventory -metrics) should be valid")
}
func TestIncludeFilters(t *testing.T) {
	testRegex, _ := regexp.Compile("four-.*")
	var testArgs = RabbitMQArguments{
		Exchanges:        []string{"one"},
		ExchangesRegexes: []*regexp.Regexp{testRegex},
		Queues:           []string{"two"},
		QueuesRegexes:    []*regexp.Regexp{testRegex},
		Vhosts:           []string{"three"},
		VhostsRegexes:    []*regexp.Regexp{testRegex},
	}
	assert.True(t, testArgs.IncludeExchange("one"))
	assert.False(t, testArgs.IncludeExchange("false"))
	assert.True(t, testArgs.IncludeExchange("four-exchange"))

	assert.True(t, testArgs.IncludeQueue("two"))
	assert.False(t, testArgs.IncludeQueue("false"))
	assert.True(t, testArgs.IncludeQueue("four-queue"))

	assert.True(t, testArgs.IncludeVhost("three"))
	assert.False(t, testArgs.IncludeVhost("false"))
	assert.True(t, testArgs.IncludeVhost("four-vhost"))
}

func TestEntityFilter(t *testing.T) {
	testArgs := RabbitMQArguments{
		Exchanges: []string{"one"},
		Queues:    []string{"two"},
		Vhosts:    []string{"three"},
	}
	assert.True(t, testArgs.IncludeEntity("one", consts.ExchangeType, "three"))
	assert.True(t, testArgs.IncludeEntity("two", consts.QueueType, "three"))
	assert.True(t, testArgs.IncludeEntity("five", consts.NodeType, "three"))
	assert.False(t, testArgs.IncludeEntity("one", consts.ExchangeType, ""))

	testArgs = RabbitMQArguments{}
	assert.True(t, testArgs.IncludeEntity("any", consts.VhostType, "any"))
}
