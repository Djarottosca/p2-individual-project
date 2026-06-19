package paymentGateway

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"p2-individual-project/service/transaction"
)

type XenditGateway struct {
	secretKey string
	client    *http.Client
}

func NewXenditGateway() *XenditGateway {
	return &XenditGateway{
		secretKey: os.Getenv("XENDIT_SECRET_KEY"),
		client:    &http.Client{},
	}
}

// CreateInvoice manggil API Xendit buat bikin invoice + payment link.
func (g *XenditGateway) CreateInvoice(externalID string, amount int64, description string) (*transaction.Invoice, error) {
	body := map[string]interface{}{
		"external_id": externalID,
		"amount":      amount,
		"description": description,
	}
	payload, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.xendit.co/v2/invoices", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Xendit pakai Basic Auth: secret key sebagai username, password kosong
	auth := base64.StdEncoding.EncodeToString([]byte(g.secretKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, errors.New("xendit error: " + string(respBody))
	}

	var parsed struct {
		ID         string `json:"id"`
		InvoiceURL string `json:"invoice_url"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	return &transaction.Invoice{
		InvoiceID:  parsed.ID,
		InvoiceURL: parsed.InvoiceURL,
	}, nil
}
