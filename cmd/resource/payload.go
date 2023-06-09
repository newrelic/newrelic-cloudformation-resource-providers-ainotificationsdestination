package resource

import (
   "fmt"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

//
// Generic, should be able to leave these as-is
//

type Payload struct {
   model  *Model
   models []interface{}
}

func (p *Payload) SetIdentifier(g *string) {
   p.model.Id = g
}

func (p *Payload) GetIdentifier() *string {
   return p.model.Id
}

func (p *Payload) GetIdentifierKey(a model.Action) string {
   return "id"
}

var emptyString = ""

func (p *Payload) GetTagIdentifier() *string {
   return &emptyString
}

func NewPayload(m *Model) *Payload {
   return &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
}

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}

func (p *Payload) GetTags() map[string]string {
   return nil
}

func (p *Payload) HasTags() bool {
   return false
}

//
// These are usually API specific, MAY BE configured per API
//

var typeName = "NewRelic::Observability::AINotificationsDestination"

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Id: &s})
}

func (p *Payload) GetGraphQLFragment() *string {
   return p.model.Destination
}

func (p *Payload) GetVariables() map[string]string {
   vars := make(map[string]string)
   if p.model.Variables != nil {
      for k, v := range p.model.Variables {
         vars[k] = v
      }
   }

   if p.model.Id != nil {
      vars["ID"] = *p.model.Id
   }

   if p.model.Destination != nil {
      vars["FRAGMENT"] = *p.model.Destination
   }

   lqf := ""
   if p.model.ListQueryFilter != nil {
      lqf = *p.model.ListQueryFilter
   }
   vars["LISTQUERYFILTER"] = lqf

   return vars
}

func (p *Payload) GetErrorKey() string {
   return "type"
}

func (p *Payload) GetResultKey(a model.Action) string {
   switch a {
   case model.Delete:
      return "ids"
   default:
      return "id"
   }
}

func (p *Payload) NeedsPropagationDelay(a model.Action) bool {
   return true
}
func (p *Payload) GetCreateMutation() string {
   return `
mutation {
  aiNotificationsCreateDestination(accountId: {{{ACCOUNTID}}}, {{{FRAGMENT}}} ){
    destination {
      id
    }
    error {
      ... on AiNotificationsConstraintsError {
        constraints {
          dependencies
          name
        }
      }
      ... on AiNotificationsDataValidationError {
        details
        fields {
          field
          message
        }
      }
      ... on AiNotificationsResponseError {
        description
        details
        type
      }
      ... on AiNotificationsSuggestionError {
        description
        details
        type
      }
    }
  }
}
`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  aiNotificationsDeleteDestination(accountId: {{{ACCOUNTID}}}, destinationId: "{{{ID}}}") {
    error {
      description
      details
      type
    }
    ids
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
  aiNotificationsUpdateDestination(accountId: {{{ACCOUNTID}}}, {{{FRAGMENT}}}, destinationId: "{{{ID}}}") {
    destination {
      id
      name
      updatedAt
      updatedBy
    }
    error {
      ... on AiNotificationsConstraintsError {
        constraints {
          dependencies
          name
        }
      }
      ... on AiNotificationsDataValidationError {
        details
        fields {
          field
          message
        }
      }
      ... on AiNotificationsResponseError {
        description
        details
        type
      }
      ... on AiNotificationsSuggestionError {
        description
        details
        type
      }
    }
  }
}
`
}

func (p *Payload) GetReadQuery() string {
   return `
{
    actor {
        account(id: {{{ACCOUNTID}}}) {
            aiNotifications {
                destinations(filters: {id: "{{{ID}}}"}) {
                    entities {
                        id
                        type
                    }
                    error {
                        description
                        details
                        type
                    }
                    nextCursor
                    totalCount
                }
            }
        }
    }
}
`
}

func (p *Payload) GetListQuery() string {
   return `
{
    actor {
        account(id: {{{ACCOUNTID}}}) {
            aiNotifications {
                destinations (cursor: "{{{NEXTCURSOR}}}"){
                    entities {
                        id
                        type
                    }
                    error {
                        description
                        details
                        type
                    }
                    nextCursor
                    totalCount
                }
            }
        }
    }
}
`
}

func (p *Payload) GetListQueryNextCursor() string {
   return ""
}
