package basecrm

type UsersService interface {
	List()
	Get()
	Self()
}

func NewUsersService(client *Client) UsersService {
	return nil
}
