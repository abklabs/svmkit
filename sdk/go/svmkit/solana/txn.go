package solana

type TxnOptions struct {
	BlockHash            *string   `pulumi:"blockHash,optional"`
	Commitment           *string   `pulumi:"commitment,optional"`
	WithComputeUnitPrice *float64  `pulumi:"withComputeUnitPrice,optional"`
	FeePayer             *string   `pulumi:"feePayer,optional" provider:"secret"`
	From                 *string   `pulumi:"from,optional" provider:"secret"`
	WithMemo             *string   `pulumi:"withMemo,optional"`
	KeyPair              *string   `pulumi:"keyPair,optional" provider:"secret"`
	Nonce                *string   `pulumi:"nonce,optional"`
	NonceAuthority       *string   `pulumi:"nonceAuthority,optional"`
	Signer               *[]string `pulumi:"signer,optional"`
	WS                   *string   `pulumi:"ws,optional"`
	URL                  *string   `pulumi:"url,optional"`
}
