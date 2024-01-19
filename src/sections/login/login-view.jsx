/* eslint-disable no-unused-vars */
/* eslint-disable import/no-extraneous-dependencies */
import { useState } from 'react';
import { useCookies } from 'react-cookie';
import 'react-toastify/dist/ReactToastify.css';
import { toast, ToastContainer } from 'react-toastify';

import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Card from '@mui/material/Card';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import LoadingButton from '@mui/lab/LoadingButton';
import { alpha, useTheme } from '@mui/material/styles';
import InputAdornment from '@mui/material/InputAdornment';

import { useRouter } from 'src/routes/hooks';
import { RouterLink } from 'src/routes/components';

import { userAPI } from 'src/api';
import { bgGradient } from 'src/theme/css';

import Iconify from 'src/components/iconify';

// ----------------------------------------------------------------------

export default function LoginView() {
  const theme = useTheme();

  const router = useRouter();

  const [cookies, setCookie] = useCookies(["token"]);

  const [input, setInput] = useState({
    username: '',
    password: '',
  });

  const [error, setError] = useState({
    username: '',
    password: '',
  });

  const [showPassword, setShowPassword] = useState(false);

  const handleInputChange = e => {
    const { name, value } = e.target;
    setInput(prev => ({
      ...prev,
      [name]: value
    }));

    setError(prev => ({
      ...prev,
      [name]: ''
    }));
  };

  const handleClick = async () => {
    const newError = {...error}

    if(!input.username) {
      newError.username = "Please enter username" 
    }

    if(!input.password) {
      newError.password = "Please enter password" 
    }

    if (newError.username || newError.password) {
      setError(newError)
      return
    }

    try {
      const result = await userAPI.signin(input.username, input.password)
      setCookie("token", result.token, { path: "/" });
      console.log(result);
      toast.success("Login successfully!")
      router.push('/')
    } catch(err) {
      console.log("Failed to signin due to ", err)
      toast.error("Something error, please try again!")
    }
  };

  const renderForm = (
    <>
      <Stack spacing={3} marginY={3}>
        <TextField 
          name="username" 
          label="Username" 
          error={!!error.username} 
          helperText={error.username}
          onChange={handleInputChange}
        />

        <TextField
          name="password"
          label="Password"
          type={showPassword ? 'text' : 'password'}
          error={!!error.password} 
          helperText={error.password}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={() => setShowPassword(!showPassword)} edge="end">
                  <Iconify icon={showPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                </IconButton>
              </InputAdornment>
            ),
          }}
          onChange={handleInputChange}
        />
      </Stack>

      <LoadingButton
        fullWidth
        size="large"
        type="submit"
        variant="contained"
        color="inherit"
        onClick={handleClick}
      >
        Login
      </LoadingButton>
    </>
  );

  return (
    <Box
      sx={{
        ...bgGradient({
          color: alpha(theme.palette.background.default, 0.9),
          imgUrl: '/assets/background/overlay_4.jpg',
        }),
        height: 1,
      }}
    >
      <Stack alignItems="center" justifyContent="center" sx={{ height: 1 }}>
        <Card
          sx={{
            p: 5,
            width: 1,
            maxWidth: 420,
          }}
        >
          <Typography variant="h4">Sign in to Scraping Tool</Typography>

          <Typography variant="body2" sx={{ mt: 2, mb: 5 }}>
            Don’t have an account?
            <Link component={RouterLink} href="/signup" variant="subtitle2" sx={{ ml: 0.5 }}>
              Get started
            </Link>
          </Typography>

          {renderForm}
        </Card>
      </Stack>
      <ToastContainer/>
    </Box>
  );
}
