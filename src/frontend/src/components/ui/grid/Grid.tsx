import React, {CSSProperties} from 'react';
import classes from './Grid.module.scss';
import {Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow} from "@material-ui/core";
import EditOutlinedIcon from '@material-ui/icons/EditOutlined';
import DeleteForeverOutlinedIcon from '@material-ui/icons/DeleteForeverOutlined';
import AddIcon from '@material-ui/icons/Add';

type Column = {
  title: string;
  field: string;
  width?: string;
}

type GridProps<d extends dataType> = {
  columns: Column[];
  data: d[];
  setEditHandler: (data: d) => void;
  addRowHandler: () => void;
  deleteHandler: (data: d) => void;
  onDoubleClick?: (data: d) => void;
}

type dataType = { id?: string; values: { key: string; values: string[] | null }[] };

function Grid<d extends dataType>(props: GridProps<d>): JSX.Element {
  const sortData = (gridData: d[]): d[] => {
    return gridData.sort((da, db) => {
      if (da.id && db.id) {
        if (da.id < db.id) {
          return -1;
        }
        if (da.id > db.id) {
          return 1;
        }
      }
      return 0;
    });
  };

  const getRow = (c: Column, gridData: d): JSX.Element => {
    const cellStyle: CSSProperties = c.width ? {maxWidth: c.width, width: c.width} : {};
    let row: JSX.Element;
    if (c.field === 'edit') {
      const {deleteHandler, setEditHandler} = props;
      row = <>
        <div style={{float: 'left'}} onClick={() => setEditHandler(gridData)}>
          <EditOutlinedIcon className={'icon'}/>
        </div>
        <div style={{float: 'left'}} onClick={() => deleteHandler(gridData)}>
          <DeleteForeverOutlinedIcon className={'icon'}/>
        </div>
      </>;
    } else {
      row = <>{gridData.values.find(value => value.key === c.field)?.values?.join(', ')}</>;
    }
    return (<TableCell
      className={classes.cell}
      style={cellStyle}
      key={c.field}
    >
      {row}
    </TableCell>);
  };

  const {addRowHandler, columns, data, onDoubleClick} = props;
  const header = columns.map(c => (<TableCell
    className={classes.head + " " + classes.cell}
    style={c.width ? {maxWidth: c.width, width: c.width} : {}}
    key={c.field}
  >
    {c.title}
  </TableCell>));
  const rows = sortData(data).map((d, i) => (<TableRow
    onDoubleClick={() => onDoubleClick ? onDoubleClick(d) : null}
    className={classes.row}
    key={d.id ?? i}
  >
    {columns.map(c => getRow(c, d))}
  </TableRow>));
  return (<>
    <div aria-label='add entry' title='add entry' onClick={addRowHandler} style={{float: 'right'}}>
      <AddIcon className={'icon'}/>
    </div>
    <TableContainer component={Paper}>
      <Table className={classes.table}>
        <TableHead>
          <TableRow>
            {header}
          </TableRow>
        </TableHead>
        <TableBody>
          {rows}
        </TableBody>
      </Table>
    </TableContainer>
  </>);
}

export default Grid;
