package dns

import "fmt"

//Represents the type of DNS Message - Request or Response.
type MessageType uint8

//Represents a DNS Message (both Request and Response).
type Message struct {
	//Represents all the data present in header section of the DNS Message
	Header Header
	//Array of Question records in the DNS Message
	Questions []Question
	//Array of answer RRs present in the DNS Message.
	Answers []Resource
	//Array of authoritative RRs present in the DNS Message.
	Authoritative []Resource
	//Array of additional RRs present in the DNS Message.
	Additional []Resource
}

//Initialises all the properties in the Message instance.
func (msg *Message) Initialize(mt MessageType) {
	msg.Header = Header{}
	msg.Header.Initialize(mt)
	msg. Questions = make([]Question, 0)
	msg.Answers = make([]Resource, 0)
	msg.Authoritative = make([]Resource, 0)
	msg.Additional = make([]Resource, 0)
}

//Set Message values for given domain name and dns record type.
func (msg *Message) NewQuestion(name string, recType RecordType) {
	question := Question{}
	question.Set(name, recType)
	msg.Questions = append(msg.Questions, question)
	questionCount := msg.Header.QdCount + 1
	msg.Header.SetQuestionCount(questionCount)
}

//Pack the message as a sequence of octets.
func (msg *Message) Pack() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, msg.Header.Pack()...)
	if msg.Header.QdCount > 0 {
		for _, que := range msg.Questions {
			buffer = append(buffer, que.Pack()...)
		}
	}
	return buffer
}

//Unpack the sequence of bytes to a Message instance.
func (msg *Message) Unpack(response []byte) {
	offset := 0
	offset = msg.Header.Unpack(response, offset)
	if msg.Header.QdCount > 0 {
		for index := 1; index <= int(msg.Header.QdCount); index++ {
			question := Question{}
			offset = question.Unpack(response, offset)
			msg.Questions = append(msg.Questions, question)
		}
	}

	if msg.Header.AnCount > 0 {
		for index := 1; index <= int(msg.Header.AnCount); index++ {
			answer := Resource{}
			offset = answer.Unpack(response, offset)
			msg.Answers = append(msg.Answers, answer)
		}
	}

	if msg.Header.NsCount > 0 {
		for index := 1; index <= int(msg.Header.AnCount); index++ {
			authoritative := Resource{}
			offset = authoritative.Unpack(response, offset)
			msg.Authoritative = append(msg.Authoritative, authoritative)
		}
	}

	if msg.Header.ArCount > 0 {
		for index := 1; index <= int(msg.Header.AnCount); index++ {
			additional := Resource{}
			offset = additional.Unpack(response, offset)
			msg.Additional = append(msg.Additional, additional)
		}
	}
}

//Returns a string representation of the DNS Message instance. 
func (msg *Message) String() string {
	string_value := fmt.Sprintf("%s\n", msg.Header.String())
	if msg.Header.QdCount > 0 {
		string_value += "QUESTION SECTION:\n"
		for _, que := range msg.Questions {
			string_value += que.String()
		}
		string_value += "\n"
	}
	
	if msg.Header.AnCount > 0 {
		string_value += "ANSWER SECTION:\n"
		for _, ans := range msg.Answers {
			string_value += ans.String()
		}
		string_value += "\n"
	}
	
	if msg.Header.NsCount > 0 {
		string_value += "AUTHORITY SECTION:\n"
		for _, auth := range msg.Authoritative {
			string_value += auth.String()
		}
		string_value += "\n"
	}

	if msg.Header.ArCount > 0 {
		string_value += "ADDITIONAL SECTION:\n"
		for _, add := range msg.Additional {
			string_value += add.String()
		}
		string_value += "\n"
	}

	return string_value
}