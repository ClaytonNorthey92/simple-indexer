import React from 'react';
import * as request from 'request-promise';
import { ListGroup } from 'react-bootstrap';

export default class Jobs extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            jobs: []
        }

        this.refreshJobs = this.refreshJobs.bind(this)
    }
    
    refreshJobs() {
        const component = this;
        request.get("http://localhost:8080/jobs")
        .then((results) => {
            component.setState({
                jobs: JSON.parse(results).map(r => ({
                    status: r.status,
                    indexedPageCount: r['indexed-page-count'],
                    indexedWordCount: r['new-words-added-count'],
                    id: r.id,
                    details: r.details,
                    startUrl: r['start-url']
                }))
            })
        }).catch(e => {
            // handle error, maybe toast?
        })
    }

    componentDidMount() {
        this.refreshJobs();
        const id = setInterval(this.refreshJobs, 2000);
        this.setState({
            intervalId: id
        })
    }

    componentWillUnmount() {
        clearInterval(this.state.intervalId)
    }

    render() {
        const statuses = {
            'in progress': 'info',
            'cancelled': "secondary",
            'succeeded': 'success',
            'failed': 'danger'
        }

        const jobsDom = this.state.jobs.map(j => (
            <ListGroup.Item key={j.id} variant={statuses[j.status]}>
                <div>
                    Index Job {j.id}
                </div>
                <div>
                    Starting URL: {j.startUrl}
                </div>
                <div>
                    Status: {j.status}
                </div>
                <div>
                    Indexed {j.indexedPageCount} Pages
                </div>
                <div>
                    Indexed {j.indexedWordCount} Words    
                </div>
                <div>
                    <b>DETAILS: {j.details}</b>
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