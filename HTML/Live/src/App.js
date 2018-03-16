import React, { Component } from 'react';
import './App.css';
import json from './test.json';
import {FileInfo} from './FileComponents';

class App extends Component {

  state = { json }
  render() {
    console.log(json);
    return (
      <div className="App">
        Test Stuff Here
        <table style={{ width: 100 }}  >
          <thead>
            <tr>
              <th>Name</th>
              <th>Size</th>
              <th>Dir</th>
            </tr>
          </thead>
          <tbody>
            <FileInfo {...this.state.json} margin={0} key={"base"} />
          </tbody>
        </table>
      </div>
    );
  }
}

export default App;
