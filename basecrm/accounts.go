package basecrm

type AccountsService interface {
	Self()
}

func NewAccountsService(client *Client) AccountsService {
	return nil
}
