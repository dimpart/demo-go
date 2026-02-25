package db

import . "github.com/dimchat/mkm-go/protocol"

const (
	AnyStation   = "station@anywhere"   // MTA: Message Transfer Agent
	AnyArchivist = "archivist@anywhere" // Profile manager
	AnyAssistant = "assistant@anywhere" // Group manager
)

type AddressNameDBI interface {

	/**
	 *  Get ID by short name
	 *
	 * @param alias - short name
	 * @return user ID
	 */
	GetIdentifier(alias string) ID

	/**
	 *  Save ANS record
	 *
	 * @param identifier - user ID
	 * @param alias - short name
	 * @return true on success
	 */
	AddRecord(identifier ID, alias string) bool

	/**
	 *  Remove ANS record
	 *
	 * @param alias - short name
	 * @return true on success
	 */
	RemoveRecord(alias string) bool
}
