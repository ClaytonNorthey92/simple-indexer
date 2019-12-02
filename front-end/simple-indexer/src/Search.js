import React from 'react';
import * as request from 'request-promise';
import {Form, Button, ListGroup, Container, Row} from 'react-bootstrap';

export default class Search extends React.Component {

    constructor(props) {
        super(props)
        this.state = {
            query: "",
            searchResults: [],
            extraDetails: ""
        }

        this.handleSubmit = this.handleSubmit.bind(this);
        this.updateQuery = this.updateQuery.bind(this);
    }

    updateQuery(event) {
        if (event.target.value === this.state.query) {
            return;
        }
        this.setState({
            query: event.target.value
        });
    }

    handleSubmit(event) {
        event.preventDefault();
        const component = this;
        request.get("http://localhost:8080/search", {
            qs: {
                q: this.state.query
            }
        }).then(result => {
            const structuredResults = JSON.parse(result);
            if (!structuredResults.length) {
                this.setState({
                    extraDetails: "no results found",
                    searchResults: []
                })
            } else {
                component.setState({
                    extraDetails: null,
                    searchResults: structuredResults
                })
            }
        }).catch(err => {
        })
    }
    
    render() {
        let extraDetails = this.state.extraDetails;
        let results = this.state.searchResults.map(s => (
            <ListGroup.Item key={s.url}>
                <div>
                    <a href={s.url} target="_blank">{s.title}</a>
                </div>
                <div>
                    Occurances: {s.count}
                </div>
            </ListGroup.Item>
        ))

        if (results && results.length) {
            extraDetails = `found ${results.length} results`;
        }

        return (
            <Container>
                <Row>
                    <Form onSubmit={this.handleSubmit}>
                        <Form.Group controlId="formIndexer">
                            <Form.Label>
                                Search Term
                            </Form.Label>
                            <Form.Control type="text" onChange={this.updateQuery}>
                            </Form.Control>
                            <Form.Text>
                                the term you want to search for
                            </Form.Text>
                        </Form.Group>
                        <Button type="submit" variant="primary">
                            Search
                        </Button>
                    </Form>
                </Row>
                <Row>
                    <div>
                        {extraDetails}
                    </div>
                    <ListGroup>
                        {results}
                    </ListGroup>
                </Row>
            </Container>
        );
    }
}