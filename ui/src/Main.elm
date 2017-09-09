port module Main exposing (main)

import Html exposing (Html, Attribute, program, text, div, h1, input)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Http exposing (Request, Response, emptyBody, expectStringResponse, request, send)
import Table
import String exposing (join)
import Json.Decode as Decode
import Json.Decode.Pipeline as Pipeline
import UUID

main : Program Never Model Msg
main =
    program
        { init = init
        , subscriptions = subscriptions
        , update = update
        , view = view
        }

-- Model

type alias Model =
    { ledgers : List Ledger
    , tableState : Table.State
    , query : Result String String
    , error : Maybe String
    }

init : ( Model, Cmd Msg )
init =
    let
        model =
         { ledgers = []
         , tableState = Table.initialSort "CreatedOn"
         , query = Ok ""
         , error = Nothing
         }
    in
        ( model, Cmd.none )

-- Subscriptions

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- Update

type Msg
    = Input String
    | SetTableState Table.State
    | RecordLedger (Result Http.Error (List Ledger))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Input (value) -> 
            let 
                q = UUID.decode value
                c = case q of
                    Ok x -> queryLedger x
                    Err _ -> Cmd.none
            in
                ( { model | query = q }, c )
        
        SetTableState (value) ->
            ( { model | tableState = value }, Cmd.none )

        RecordLedger (Ok ledgers) ->
            ( { model | ledgers = ledgers }, Cmd.none )
        
        RecordLedger (Err err) -> 
            let
                m = case err of
                    Http.BadStatus response -> response.body
                    Http.BadPayload debug response -> debug ++ response.body
                    _ -> toString err
            in
                ( { model | error = Just m }, Cmd.none )

-- View

view : Model -> Html Msg
view { ledgers, tableState, query, error } =
    let
        instructionsText = case query of
            Ok _ -> case error of
                Just msg -> ( errorViewStyle, msg )
                Nothing -> ( instructionsViewStyle, "" )
            Err x -> ( errorViewStyle, x )
    in
        div []
            [ h1 [ h1ViewStyle ] [ text "Ledgers" ]
            , input [ placeholder "Search by ResourceID", onInput Input, inputViewStyle ] []
            , div [ Tuple.first instructionsText ] [ text <| Tuple.second instructionsText ]
            , div [ tableViewStyle ] 
                [ Table.view config tableState ledgers 
                ]
            ]

config : Table.Config Ledger Msg
config =
    Table.config
        { toId = .name
        , toMsg = SetTableState
        , columns =
            [ Table.stringColumn  "Name" .name 
            , Table.stringColumn  "ParentID" .parent_id
            , Table.stringColumn  "ResourceID" .resource_id
            , Table.stringColumn  "ResourceAddress" .resource_address
            , Table.intColumn     "ResourceSize" .resource_size
            , Table.stringColumn  "ResourceContentType" .resource_content_type
            , Table.stringColumn  "AuthorID" .author_id
            , Table.stringColumn  "Tags" <| join "," << .tags 
            , Table.stringColumn  "CreatedOn" .created_on
            , Table.stringColumn  "DeletedOn" .deleted_on
            ]
        }

h1ViewStyle : Attribute Msg
h1ViewStyle = 
    style
        [ ("width", "100%")
        , ("height", "40px")
        , ("padding", "10px 0")
        , ("font-size", "2em")
        , ("text-align", "center")
        , ("background-color", "#ffffff")
        , ("color", "#4ac64a")
        ]

tableViewStyle : Attribute Msg
tableViewStyle =
    style
        [ ("margin", "10px 20px")
        ]

inputViewStyle : Attribute Msg
inputViewStyle =
    style
        [ ("font-size", "2em")
        , ("margin", "10px 20px")
        , ("color", "#4ac64a")
        ]

instructionsViewStyle : Attribute Msg
instructionsViewStyle =
    style
        [ ("display", "none") 
        ]

errorViewStyle : Attribute Msg
errorViewStyle =
    style
        [ ("display", "block")
        , ("margin", "10px 0")
        , ("width", "100%")
        , ("height", "40px")
        , ("padding", "10px 0")
        , ("font-size", "2em")
        , ("text-align", "center")
        , ("background-color", "#f44336")
        , ("color", "#ffffff")
        ]

-- Ledger

type alias Ledger =
    { name : String
    , parent_id : String
    , resource_id : String
    , resource_address : String
    , resource_size : Int
    , resource_content_type : String
    , author_id : String
    , tags : List String
    , created_on : String
    , deleted_on : String
    }

-- API

queryLedger : String -> Cmd Msg
queryLedger resource =
    let
        url =
            "/ledgers/multiple/?resource_id=" ++ resource
    in
        Http.send RecordLedger ( Http.get url decodeLedgers )
        
decodeLedgers : Decode.Decoder (List Ledger)
decodeLedgers =
    Decode.list decodeLedger

decodeLedger : Decode.Decoder Ledger
decodeLedger =
    Pipeline.decode Ledger
        |> Pipeline.required "name" Decode.string
        |> Pipeline.required "parent_id" Decode.string
        |> Pipeline.required "resource_id" Decode.string
        |> Pipeline.required "resource_address" Decode.string
        |> Pipeline.required "resource_size" Decode.int
        |> Pipeline.required "resource_content_type" Decode.string
        |> Pipeline.required "author_id" Decode.string
        |> Pipeline.required "tags"  (Decode.list Decode.string)
        |> Pipeline.required "created_on" Decode.string
        |> Pipeline.required "deleted_on" Decode.string