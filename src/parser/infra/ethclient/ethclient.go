package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/apm-dev/evm-tx-parser/src/common"
	"github.com/apm-dev/evm-tx-parser/src/config"
	"github.com/apm-dev/evm-tx-parser/src/domain"
	"github.com/pkg/errors"
)

const (
	DefaultRequestID = int64(1)
	JsonRpcVersion   = "2.0"
)

type ethClient struct {
	config *config.Config
	url    string
}

func NewEthClient(
	config *config.Config,
	url string,
) domain.EthereumClient {
	return &ethClient{
		config: config,
		url:    url,
	}
}

func (c *ethClient) GetNowBlockNumber() (int, error) {
	rsp := &RpcResponse{}
	err := c.sendRequest(&RpcRequest{
		ID:      DefaultRequestID,
		JsonRpc: JsonRpcVersion,
		Method:  "eth_blockNumber",
		Params:  nil,
	}, rsp)
	if err != nil {
		return 0, err
	}
	if rsp.Error != nil {
		return 0, errors.Errorf("eth_blockNumber failed, code: %d, msg: %s, data: %s", rsp.Error.Code, rsp.Error.Message, rsp.Error.Data)
	}
	var num string
	err = json.Unmarshal(rsp.Result, &num)
	if err != nil {
		return 0, err
	}
	return int(common.HexToInt(num)), nil
}

func (c *ethClient) GetBlocksByRange(from, to int) ([]domain.Block, error) {
	var rsps []RpcResponse
	reqs := make([]RpcRequest, 0, to-from)
	for i := from; i < to+1; i++ {
		reqs = append(reqs, RpcRequest{
			ID:      int64(i),
			JsonRpc: JsonRpcVersion,
			Method:  "eth_getBlockByNumber",
			Params: []any{
				common.IntToHex(i),
				true,
			},
		})
	}
	err := c.sendRequest(reqs, &rsps)
	if err != nil {
		return nil, err
	}
	if len(rsps) == 0 {
		return nil, nil
	}
	blocks := make([]domain.Block, 0)
	for _, rsp := range rsps {
		if rsp.Error != nil {
			return nil, errors.Errorf("eth_getBlockByNumber failed, code: %d, msg: %s, data: %s", rsp.Error.Code, rsp.Error.Message, rsp.Error.Data)
		}
		var block Block
		err = json.Unmarshal(rsp.Result, &block)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, *block.ToEntity())
	}
	return blocks, nil
}

func (c *ethClient) sendRequest(body any, result any) error {
	const op = "ethclient.sendRequest"

	reqBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(domain.ErrInternalServer, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.config.App.OutgoingRequestTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.Wrap(domain.ErrInternalServer, err.Error())
	}
	request.Header.Add("Content-Type", "application/json")
	// send request to node and parse response
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(domain.ErrInternalServer, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		rsp, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(domain.ErrInternalServer, "%s: request to '%s' failed with status: '%d'", op, c.url, resp.StatusCode)
		} else {
			return errors.Wrapf(domain.ErrInternalServer, "%s: request to '%s' failed with status: '%d', and message '%s'", op, c.url, resp.StatusCode, string(rsp))
		}
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(domain.ErrInternalServer, err.Error())
	}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return errors.Wrap(domain.ErrInternalServer, err.Error())
	}
	return nil
}
