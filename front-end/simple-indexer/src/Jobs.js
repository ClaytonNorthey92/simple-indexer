import React from 'react';
import * as request from 'request-promise';
import { ListGroup } from 'react-bootstrap';

export default class Jobs extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            jobs: []
        }
    }
    
    componentDidMount() {
        const component = this;
        setInterval(() => {
            request.get("http://localhost:8080/jobs")
            .then((results) => {
                component.setState({
                    jobs: JSON.parse(results).map(r => ({
                        status: r.status,
                        indexedPageCount: r['indexed-page-count'],
                        indexedWordCount: r['indexed-word-count'],
                        id: r.id,
                        details: r.details
                    }))
                })
            }).catch(e => {
                // handle error, maybe toast?
            })
        }, 500)

    }

    render() {
        const statuses = {
            'in progress': 'info',
            'succeeded': 'success',
            'failed': 'danger'
        }

        const jobsDom = this.state.jobs.map(j => (
            <ListGroup.Item key={j.id} variant={statuses[j.status]}>
                <div>
                    Index Job {j.id}
                </div>
                <div>
                    {j.status}
                </div>
                <div>
                    {j.details}
                </div>
            </ListGroup.Item>

        ))
        return (
            <ListGroup>
                    {jobsDom}
            </ListGroup>
        );
    }
}