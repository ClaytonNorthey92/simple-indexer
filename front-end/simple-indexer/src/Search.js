import React from 'react';
import * as request from 'request-promise';
import {Form, Button, ListGroup, Container, Row} from 'react-bootstrap';

export default class Search extends React.Component {

    constructor(props) {
        super(props)
        this.state = {
            query: "",
            searchResults: []
        }

        this.handleSubmit = this.handleSubmit.bind(this);
        this.updateQuery = this.updateQuery.bind(this);
    }

    updateQuery(event) {
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
            component.setState({
                searchResults: JSON.parse(result)
            })
        }).catch(err => {
            debugger;
        })
    }
    
    render() {
        const results = this.state.searchResults.map(s => (
            <ListGroup.Item key={s.url}>
                <div>
                    <a href={s.url} >{s.title}</a>
                </div>
                <div>
                    Found {s.count} matches for "{this.state.query}"
                </div>
            </ListGroup.Item>
        ))

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
                    <ListGroup>
                        {results}
                    </ListGroup>
                </Row>
            </Container>
        );
    }
}