import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  TextField, 
  Button, 
  Paper, 
  CircularProgress,
  Stack,
} from '@mui/material';
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

function AIQuestsTab() {
  const [prompt, setPrompt] = useState('');
  const [loading, setLoading] = useState(false);

  const generateAIQuest = async (e) => {
    e.preventDefault();
    
    if (!prompt.trim()) {
      showToast('Введите описание квеста для генерации', 'warning');
      return;
    }

    setLoading(true);
    
    try {
      const result = await apiCall('/quests/generate', {
        method: 'POST',
        body: JSON.stringify({ prompt: prompt.trim() }),
      });

      if (result.success) {
        showToast('Квест успешно сгенерирован и добавлен в магазин!', 'success');
        setPrompt('');
      } else {
        showToast('Ошибка генерации квеста: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка при обращении к серверу: ' + error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 900, mx: 'auto' }}>
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 4 }}>
        <AutoAwesomeIcon sx={{ fontSize: 48, color: 'primary.main' }} />
        <Box>
          <Typography variant="h3" gutterBottom sx={{ fontWeight: 700 }}>
            Создать квест с помощью AI
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Опишите желаемый квест, и AI создаст его для вас
          </Typography>
        </Box>
      </Stack>

      <Paper 
        elevation={0}
        sx={{ 
          p: 4,
          border: '1px solid',
          borderColor: 'divider',
        }}
      >
        <form onSubmit={generateAIQuest}>
          <Stack spacing={3}>
            <TextField
              fullWidth
              multiline
              rows={10}
              placeholder="Опишите, какой квест вы хотите создать.

Например: 'Квест для развития силы воли и физического здоровья'
Или: 'Интеллектуальный квест для улучшения памяти'
Или: 'Социальный квест для развития коммуникативных навыков'"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              disabled={loading}
              label="Описание квеста"
              helperText="Будьте максимально конкретны в описании желаемого квеста"
            />

            {loading && (
              <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 3 }}>
                <CircularProgress size={50} />
                <Typography variant="body1" color="text.secondary" sx={{ mt: 2 }}>
                  ИИ генерирует ваш квест...
                </Typography>
              </Box>
            )}

            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={loading}
              size="large"
              startIcon={loading ? null : <AutoAwesomeIcon />}
            >
              {loading ? 'Генерация...' : 'Сгенерировать квест'}
            </Button>
          </Stack>
        </form>
      </Paper>
    </Box>
  );
}

export default AIQuestsTab;
