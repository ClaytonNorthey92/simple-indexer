import React from 'react';
import './App.css';
import Jobs from './Jobs';
import Search from './Search';
import Indexer from './Indexer';
import { Row, Container, Col, ButtonToolbar, Button } from 'react-bootstrap'
import * as request from 'request-promise';

class App extends React.Component {
  activeScreenIndexer = "indexer"
  activeScreenSearch = "search"

  constructor(props) {
    super(props)
    this.state = {
      activeScreen: this.activeScreenIndexer
    }
    this.setActiveScreenIndexer = this.setActiveScreenIndexer.bind(this)
    this.setActiveScreenSearch  = this.setActiveScreenSearch.bind(this)
  }

  deleteIndex() {
    request.delete("http://localhost:8080/index")
  }

  setActiveScreenIndexer() {
    this.setState({
      activeScreen: this.activeScreenIndexer
    })
  }

  setActiveScreenSearch() {
    this.setState({
      activeScreen: this.activeScreenSearch
    })
  }

  render() {
    let body;
    if (this.state.activeScreen === this.activeScreenIndexer) {
      body = (
        <Row>
          <Col>
            <Indexer></Indexer>
          </Col>
          <Col>
            <Jobs></Jobs>
          </Col>
        </Row>
      );
    }

    if (this.state.activeScreen === this.activeScreenSearch) {
      body = (
        <Row>
          <Col>
            <Search></Search>
          </Col>
        </Row>
      );
    }

    return (
      <div className="App">
        <Container>
          <Row>
            <ButtonToolbar>
              <Col>
                <Button variant={this.state.activeScreen === this.activeScreenIndexer ? "primary" : "outline-primary"} onClick={this.setActiveScreenIndexer}>
                  Index
                </Button>
              </Col>
              <Col>
                <Button variant={this.state.activeScreen === this.activeScreenSearch ? "primary" : "outline-primary"} onClick={this.setActiveScreenSearch}>
                  Search
                </Button>
              </Col>
            </ButtonToolbar>
          </Row>
          {body}
        </Container>
      </div>
    );
  }
}

export default App;
