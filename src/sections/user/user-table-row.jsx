import PropTypes from 'prop-types';

import Button from '@mui/material/Button';
import TableRow from '@mui/material/TableRow';
import TableCell from '@mui/material/TableCell';

// ----------------------------------------------------------------------

export default function UserTableRow({
  index,
  keyword,
  createdAt,
  linkTotal,
  adwordTotal,
  searchResultTotal,
  handleHtmlContentClick
}) {
  return (
    <TableRow hover tabIndex={-1}>
      <TableCell component="th" scope="row" padding="normal" align="center">
        {index}
      </TableCell>
      <TableCell>{keyword}</TableCell>
      <TableCell>{createdAt}</TableCell>
      <TableCell>{linkTotal}</TableCell>
      <TableCell>{adwordTotal}</TableCell>
      <TableCell>{searchResultTotal}</TableCell>
      <TableCell align="center">
        <Button onClick={handleHtmlContentClick} variant="contained" type="primary">
          Open
        </Button>
      </TableCell>
    </TableRow>
  );
}

UserTableRow.propTypes = {
  index: PropTypes.number,
  keyword: PropTypes.string,
  createdAt: PropTypes.any,
  linkTotal: PropTypes.number,
  adwordTotal: PropTypes.number,
  searchResultTotal:PropTypes.any,
  handleHtmlContentClick: PropTypes.func,
};
