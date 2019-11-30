import React from 'react';
import './App.css';
import Jobs from './Jobs';
import Indexer from './Indexer';
import { Row, Container, Col } from 'react-bootstrap'

function App() {
  return (
    <div className="App">
      <Container>
        <Row>
        <Col>
          <Indexer></Indexer>
        </Col>
        <Col>
          <Jobs></Jobs>
        </Col>
        </Row>
      </Container>
    </div>
  );
}

export default App;
