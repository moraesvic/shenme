package types

type ILength interface {
	Length() int
}

type IDefinition interface {
	ILength
}

type IDefinitionString interface {
	ILength
	String() string
}

type IDefinitionHTML interface {
	IDefinition
	HTML() string
}

type IDefinitionBoth interface {
	IDefinitionString
	IDefinitionHTML
}

type IDefinerString interface {
	Define(traditional string) IDefinitionString
}

type IDefinerHTML interface {
	DefineHTML(traditional string) IDefinitionHTML
}

type IDefinerBoth interface {
	IDefinerString
	IDefinerHTML
}
