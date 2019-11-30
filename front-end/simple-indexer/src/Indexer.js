import React from 'react'
import { Form, Button } from 'react-bootstrap';
import * as request from 'request-promise';

export default class Indexer extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            url: "",
            urlInvalid: false
        }
        this.handleSubmit = this.handleSubmit.bind(this);
        this.updateQuery = this.updateQuery.bind(this);
    }


    handleSubmit(event) {
        event.preventDefault();
        const component = this;
        request({
            method: "POST",
            uri: "http://localhost:8080/index",
            body: {
                url: this.state.url
            },
            json: true
        })
        .then((results) => {
            component.setState({
                urlInvalid: false
            })
        }).catch(err => {
            component.setState({
                urlInvalid: true
            })
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
                <Form.Control isInvalid={this.state.urlInvalid} type="text" onChange={this.updateQuery}>
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