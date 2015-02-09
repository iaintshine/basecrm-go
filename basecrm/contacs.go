package basecrm

type ContactsService interface {
	List()
	Get()
	Create()
	Edit()
	Delete()
}

func NewContactsService(client *Client) ContactsService {
	return nil
}
