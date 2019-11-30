import React from 'react'
import { Form, Button } from 'react-bootstrap';
import * as request from 'request-promise';

export default class Indexer extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            url: ""
        }
        this.handleSubmit = this.handleSubmit.bind(this);
        this.updateQuery = this.updateQuery.bind(this);
    }


    handleSubmit(event) {
        event.preventDefault();
        request({
            method: "POST",
            uri: "http://localhost:8080/index",
            body: {
                url: this.state.url
            },
            json: true
        })
        .then((results) => {
        }).catch(err => {
        })
    }

    updateQuery(event) {
        const url = event.target.value;
        this.setState({
            url: url
        })
    }

    render() {
        return (
            <Form onSubmit={this.handleSubmit}>
            <Form.Group controlId="formIndexer">
                <Form.Label>
                    URL
                </Form.Label>
                <Form.Control type="text" onChange={this.updateQuery}>
                </Form.Control>
                <Form.Text>
                    enter the url you want to start indexing from
                </Form.Text>
            </Form.Group>
            <Button type="submit" variant="primary">
                Index
            </Button>
        </Form>
        )
    }
}