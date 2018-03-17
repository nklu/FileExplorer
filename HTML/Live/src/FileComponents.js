import React, { Component } from 'react';
import './App.css';

const PlusMinusButton = ({ open, onClick }) =>
  (
    <button value={open} onClick={onClick}>
      {open ? "-" : "+"}
    </button>
  );

const DisplayName = ({ name }) => <div>{name}</div>
const DisplaySize = ({ size }) => <div>{bytesToSize(size)}</div>
const DisplaySizePercentage = ({ size, parentSize }) => <div>{((size / parentSize) * 100).toFixed(2)}%</div>
const DisplayDirectory = ({ isDir }) => <div>{isDir ? "True" : "False"}</div>
const DisplayExpandibleName = ({ name, open, onClick, hasChildren, margin }) => (
  <div style={{ marginLeft: margin, display: "flex" }}>
    {hasChildren && <PlusMinusButton open={open} onClick={onClick} />}
    <DisplayName name={name} />
  </div>
);

const DisplayData = ({ data, hasChildren, open, onClick, margin, parentSize }) => {
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
          <DisplaySizePercentage size={data.Size} parentSize={parentSize} />
        </td>
        <td>
          <DisplayDirectory isDir={data.IsDir} />
        </td>
      </tr>
    );
  }
}


class FileInfo extends Component {

  state = { open: false }

  handleButtonClick() {
    const { open } = this.state;
    this.setState({ ...this.state, open: !open })
  }

  render() {
    const { Data, Children, margin, parentSize, ...mainData } = this.props;
    const { open } = this.state;
    const hasChildren = Children && Children.length > 0;
    const sortedChildren = (Children && Children.sort((a, b) => b.Size - a.Size)) || [];
    const topParentSize = parentSize || this.props.Size;

    return [
      <DisplayData
        key={1}
        data={mainData}
        open={open}
        onClick={this.handleButtonClick.bind(this)}
        margin={margin}
        hasChildren={hasChildren}
        parentSize={topParentSize}
      />,
      open &&
      sortedChildren.map(child =>
        <FileInfo key={child.Name} {...child} margin={margin + 10} parentSize={this.props.Size} />
      )
    ];
  }
}

const FileInfoHeader = () => (
  <tr>
    <th>Name</th>
    <th>Size</th>
    <th>Size Percent</th>
    <th>Dir</th>
  </tr>
);

export { FileInfo, FileInfoHeader };

// https://stackoverflow.com/questions/15900485/correct-way-to-convert-size-in-bytes-to-kb-mb-gb-in-javascript
const bytesToSize = (bytes) => {
  var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 Byte';
  var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)), 10);
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + sizes[i];
};