import React, { useState } from 'react';
import { 
  Box, 
  TextField, 
  Button, 
  Paper, 
  Typography, 
  Stack,
  InputAdornment,
  IconButton,
} from '@mui/material';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import PersonIcon from '@mui/icons-material/Person';
import EmailIcon from '@mui/icons-material/Email';
import LockIcon from '@mui/icons-material/Lock';
import { login, register } from '../api';
import { showToast } from '../utils/toast';

function AuthTab({ onLogin }) {
  const [isRegistering, setIsRegistering] = useState(false);
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!username || !password) {
      showToast('Заполните все поля', 'warning');
      return;
    }

    setLoading(true);
    try {
      const result = await login(username, password);
      if (result.success) {
        showToast('Успешный вход! Добро пожаловать!', 'success');
        setTimeout(() => onLogin(), 500);
      } else {
        showToast(result.error || result.data?.error || 'Ошибка входа', 'error');
      }
    } catch (error) {
      showToast('Ошибка входа: ' + error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    if (!username || !password) {
      showToast('Заполните все обязательные поля', 'warning');
      return;
    }

    setLoading(true);
    try {
      const result = await register(username, email, password);
      if (result.success) {
        showToast('Регистрация успешна! Теперь войдите.', 'success');
        setIsRegistering(false);
        setEmail('');
        setPassword('');
      } else {
        showToast(result.data?.error || 'Ошибка регистрации', 'error');
      }
    } catch (error) {
      showToast('Ошибка регистрации: ' + error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 500, mx: 'auto', py: 8 }}>
      <Typography 
        variant="h3" 
        align="center" 
        gutterBottom 
        sx={{ 
          mb: 6, 
          fontWeight: 700,
        }}
      >
        {isRegistering ? 'Регистрация' : 'Вход'}
      </Typography>
      
      <Paper 
        elevation={0}
        sx={{ 
          p: 4,
          border: '1px solid',
          borderColor: 'divider',
        }}
      >
        <form onSubmit={isRegistering ? handleRegister : handleLogin}>
          <Stack spacing={3}>
            <TextField
              fullWidth
              label="Имя пользователя"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <PersonIcon />
                  </InputAdornment>
                ),
              }}
            />
            
            {isRegistering && (
              <TextField
                fullWidth
                type="email"
                label="Email (опционально)"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <EmailIcon />
                    </InputAdornment>
                  ),
                }}
              />
            )}
            
            <TextField
              fullWidth
              type={showPassword ? 'text' : 'password'}
              label="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <LockIcon />
                  </InputAdornment>
                ),
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowPassword(!showPassword)}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            <Button
              type="submit"
              variant="contained"
              fullWidth
              size="large"
              disabled={loading}
              sx={{ mt: 2 }}
            >
              {loading ? 'Загрузка...' : (isRegistering ? 'Зарегистрироваться' : 'Войти')}
            </Button>
          </Stack>
        </form>

        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Button
            fullWidth
            variant="text"
            onClick={() => {
              setIsRegistering(!isRegistering);
              setPassword('');
            }}
          >
            {isRegistering ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
          </Button>
        </Box>
      </Paper>
    </Box>
  );
}

export default AuthTab;
