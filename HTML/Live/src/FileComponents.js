import React, { Component } from 'react';
import './App.css';

const PlusMinusButton = ({ open, onClick }) =>
  (
    <button value={open} onClick={onClick}>
      {open ? "-" : "+"}
    </button>
  );

const DisplayName = ({ name }) => (<div>{name}</div>)
const DisplaySize = ({ size }) => (<div>{size}</div>)
const DisplayDirectory = ({ isDir }) => (<div>{isDir ? "True" : "False"}</div>)
const DisplayExpandibleName = ({ name, open, onClick, hasChildren, margin }) => (
  <div style={{ marginLeft: margin, display: "flex" }}>
    {hasChildren && <PlusMinusButton open={open} onClick={onClick} />}
    <DisplayName name={name} />
  </div>
);

const DisplayData = ({ data, hasChildren, open, onClick, margin }) => {
  if (data) {
    return (
      <tr  >
        <td>
          <DisplayExpandibleName
            name={data.Name}
            hasChildren={hasChildren}
            open={open}
            onClick={onClick}
            margin={margin}
          />
        </td>
        <td>
          <DisplaySize size={data.Size} />
        </td>
        <td>
          <DisplayDirectory isDir={data.IsDir} />
        </td>
      </tr>
    );
  }
}


export class FileInfo extends Component {

  state = { open: false }

  handleButtonClick() {
    const { open } = this.state;
    this.setState({ ...this.state, open: !open })
  }

  render() {
    const { Data, Children, margin } = this.props;
    const { open } = this.state;
    const hasChildren = Children && Children.length > 0;
    return [
      <DisplayData
        key={1}
        data={Data}
        open={open}
        onClick={this.handleButtonClick.bind(this)}
        margin={margin}
        hasChildren={hasChildren} />,
      open &&
      Children &&
      Children.map(child => <FileInfo key={child.Data.Name} {...child} margin={margin + 10} />)
    ];
  }
}

