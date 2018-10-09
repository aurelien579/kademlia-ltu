package kademlia

type ContactResult struct {
	ID      string
	Address string
}

func NewContactResult(contact *Contact) ContactResult {
	return ContactResult{
		ID:      contact.ID.String(),
		Address: contact.Address,
	}
}

func ContactsToContactResults(contacts []Contact) []ContactResult {
	var contactsResults []ContactResult

	for _, c := range contacts {
		contactsResults = append(contactsResults, NewContactResult(&c))
	}

	return contactsResults
}

func ContactResultsToContacts(results []ContactResult) []Contact {
	var contacts []Contact

	for _, r := range results {
		contacts = append(contacts, NewContact(NewKademliaID(r.ID), r.Address))
	}

	return contacts
}
