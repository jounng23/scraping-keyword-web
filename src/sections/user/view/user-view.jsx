/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable react/jsx-no-useless-fragment */
/* eslint-disable import/no-extraneous-dependencies */
import 'react-toastify/dist/ReactToastify.css';
import { useRef, useState, useEffect } from 'react';
import { toast, ToastContainer } from 'react-toastify';

import Card from '@mui/material/Card';
import Stack from '@mui/material/Stack';
import Table from '@mui/material/Table';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import TableBody from '@mui/material/TableBody';
import Typography from '@mui/material/Typography';
import TableContainer from '@mui/material/TableContainer';
import TablePagination from '@mui/material/TablePagination';

import { useRouter } from 'src/routes/hooks';

import { userAPI, keywordAPI } from 'src/api';

import Iconify from 'src/components/iconify';
import Scrollbar from 'src/components/scrollbar';

import { emptyRows } from '../utils';
import TableNoData from '../table-no-data';
import UserTableRow from '../user-table-row';
import UserTableHead from '../user-table-head';
import TableEmptyRows from '../table-empty-rows';

const defaultPage = 0;
const defaultRowsPerPage = 5;
const defaultOrder = "desc";
const defaultOrderBy = "created_at";

export default function UserPage() {
  const router = useRouter();

  const [keywordResults, setKeywordResults] = useState([]);

  const [total, setTotal] = useState(0);

  const [page, setPage] = useState(defaultPage);

  const [order, setOrder] = useState(defaultOrder);

  const [orderBy, setOrderBy] = useState(defaultOrderBy);

  const [rowsPerPage, setRowsPerPage] = useState(defaultRowsPerPage);

  const fileInputRef = useRef(null);

  const [isVerifiedUser, setIsVerifiedUser] = useState(false);

  const verifyUser = async () => {
    try {
      const result = await userAPI.verify()
      setIsVerifiedUser(true);
      console.log(result);
    } catch (err) {
      console.log(`Failed to verify user due to: ${err}`);
      toast.error("The login session has expired!");
      setTimeout(() => router.push("/login"), 5000);
    }
  }

  const checkIsCollectKeywordResultInvalid = (result) => (!result || !Array.isArray(result.keyword_results) || !result.metadata || !result.metadata.total)

  const initKeywordData = async () => {
    setPage(defaultPage);
    setRowsPerPage(defaultRowsPerPage);
    setOrder(defaultOrder);
    setOrderBy(defaultOrderBy);

    const result = await keywordAPI.collectKeywords(defaultPage + 1, defaultRowsPerPage, `${defaultOrderBy} ${defaultOrder}`)
    console.log(result);
    if(checkIsCollectKeywordResultInvalid(result)) {
      return
    }
    
    setKeywordResults(result.keyword_results);
    setTotal(result.metadata.total);
  }

  const collectKeyword = async () => {
    const result = await keywordAPI.collectKeywords(page + 1, rowsPerPage, `${orderBy} ${order}`)
    console.log(result);
    if(checkIsCollectKeywordResultInvalid(result)) {
      return
    }
    
    setKeywordResults(result.keyword_results);
    setTotal(result.metadata.total);
  }
  
  useEffect(() => {
    verifyUser()
  }, [])

  useEffect(() => {
    if(!isVerifiedUser) {
      return
    }
    initKeywordData()
  }, [isVerifiedUser])

  useEffect(() => {
    collectKeyword()
  }, [order, orderBy])

  const handleSort = (event, id) => {
    if(id !== 'created_at') {
      return
    }
    
    const isAsc = orderBy === id && order === 'asc';
    if (id !== '') {
      setOrder(isAsc ? 'desc' : 'asc');
      setOrderBy(id);
    }
  };

  const handleHtmlPageClick = (htmlCtn) => {
    const newTab = window.open();
    newTab.document.open();
    newTab.document.write(htmlCtn);
    newTab.document.close();
  } 

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setPage(0);
    setRowsPerPage(parseInt(event.target.value, 10));
  };

  const handleUploadFileChange = async (e) => {
    if(!e.target) {
      return
    }

    const files = e.target.files || []
    const file = files[0]

    if (typeof file === 'undefined' ) return;
    const formData = new FormData();
    formData.append('file', file);

    const result = await keywordAPI.uploadKeywords(formData)
    console.log(result)
    toast.success("Upload keywords successfully!")
    initKeywordData()
  }

  const handleUploadFileClick = () => {
    fileInputRef.current.value = '';
  }

  return <Container>
    <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
      <Typography variant="h4">Keyword</Typography>

      <Button component="label" variant="contained" color="success" startIcon={<Iconify icon="mdi:upload" />}>
        Upload Keywords
        <input 
          hidden 
          type="file" 
          accept=".csv" 
          ref={fileInputRef} 
          onClick={handleUploadFileClick} 
          onChange={handleUploadFileChange} 
        />
      </Button>
    </Stack>

    <Card>
      <Scrollbar>
        <TableContainer sx={{ overflow: 'unset' }}>
          <Table sx={{ minWidth: 800 }}>
            <UserTableHead
              order={order}
              orderBy={orderBy}
              rowCount={keywordResults.length}
              onRequestSort={handleSort}
              headLabel={[
                { id: 'index', label: 'Order' },
                { id: 'keyword', label: 'Keyword' },
                { id: 'created_at', label: 'Created At' },
                { id: 'link_total', label: 'Link Total' },
                { id: 'adword_total', label: 'Adword Total' },
                { id: 'search_result_total', label: 'Search Result Total' },
                { id: 'html_content', label: 'Html Content', align: 'center' },
              ]}
            />
            <TableBody>
              {keywordResults
                .map((row, i) => (
                  <UserTableRow
                    key={row.ID}
                    index={i+1}
                    keyword={row.keyword}
                    linkTotal={row.link_total}
                    adwordTotal={row.adword_total}
                    createdAt={new Date(row.created_at).toLocaleString()}
                    searchResultTotal={row.search_result_total}
                    handleHtmlContentClick={() => handleHtmlPageClick(row.html_content)}
                  />
                ))}

              <TableEmptyRows
                height={77}
                emptyRows={emptyRows(page, rowsPerPage, total)}
              />

              {(!isVerifiedUser || keywordResults.length === 0 ||  total === 0) && <TableNoData />}
            </TableBody>
          </Table>
        </TableContainer>
      </Scrollbar>

      <TablePagination
        page={page}
        component="div"
        count={total}
        rowsPerPage={rowsPerPage}
        onPageChange={handleChangePage}
        rowsPerPageOptions={[5, 10, 25]}
        onRowsPerPageChange={handleChangeRowsPerPage}
      />
    </Card>
    <ToastContainer/>
  </Container>
}
