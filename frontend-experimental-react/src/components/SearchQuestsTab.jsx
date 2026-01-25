import React, { useState } from 'react';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  InputAdornment,
  Stack,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Chip,
  Divider,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import StarIcon from '@mui/icons-material/Star';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import { apiCall } from '../api';
import QuestCard from './QuestCard';
import MasonryGrid from './MasonryGrid';
import { showToast } from '../utils/toast';

function SearchQuestsTab() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [questDetails, setQuestDetails] = useState(null);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);

  const submitQuestSearch = async (e) => {
    e.preventDefault();

    if (!query.trim()) {
      showToast('Введите текст запроса для поиска квестов.', 'warning');
      return;
    }

    setLoading(true);
    setResults(null);

    try {
      const body = {
        query: query.trim(),
        status: 'all',
        top_k: 10,
      };

      const result = await apiCall('/quests/search', {
        method: 'POST',
        body: JSON.stringify(body),
      });

      setLoading(false);

      if (!result.success) {
        showToast('Ошибка поиска: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
        return;
      }

      if (!Array.isArray(result.data)) {
        console.error('Search returned non-array data:', result.data);
        showToast('Сервис поиска вернул данные в неожиданном формате.', 'error');
        return;
      }

      // Логируем для отладки
      console.log('Search results:', result.data);
      
      setResults(result.data);
      if (result.data.length === 0) {
        showToast('По вашему запросу ничего не найдено', 'info');
      }
    } catch (error) {
      setLoading(false);
      showToast('Ошибка поиска: ' + error.message, 'error');
    }
  };

  const purchaseQuest = async (quest) => {
    // Извлекаем quest из структуры {quest: {...}, similarity_score: ...} если нужно
    const questData = quest.quest || quest;
    if (!questData || !questData.id) return;
    
    try {
      const result = await apiCall(`/quests/${questData.id}/purchase`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Квест успешно куплен!', 'success');
      } else {
        showToast('Ошибка покупки: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка покупки квеста', 'error');
    }
  };

  const showQuestDetails = async (quest) => {
    // Извлекаем quest из структуры {quest: {...}, similarity_score: ...} если нужно
    const questData = quest.quest || quest;
    if (!questData || !questData.id) return;
    
    try {
      const result = await apiCall(`/quests/${questData.id}/details`);
      if (result.success) {
        setQuestDetails(result.data);
        setDetailsDialogOpen(true);
      } else {
        showToast('Ошибка загрузки деталей', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки деталей', 'error');
    }
  };

  return (
    <Box>
      <Typography variant="h3" gutterBottom sx={{ mb: 4, fontWeight: 700 }}>
        Поиск квестов
      </Typography>

      <Paper 
        elevation={0}
        sx={{ 
          p: 4, 
          mb: 4,
          border: '1px solid',
          borderColor: 'divider',
        }}
      >
        <form onSubmit={submitQuestSearch}>
          <Stack spacing={2}>
            <TextField
              fullWidth
              placeholder="Например: 'спорт челлендж', 'развитие силы воли', 'совместный квест с другом'"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
            <Button 
              type="submit" 
              variant="contained" 
              disabled={loading}
              size="large"
              startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SearchIcon />}
            >
              {loading ? 'Поиск...' : 'Искать'}
            </Button>
            <Typography variant="body2" color="text.secondary">
              Опишите желаемый квест естественным языком — AI подберёт наиболее подходящие варианты.
            </Typography>
          </Stack>
        </form>
      </Paper>

      {loading && (
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 6 }}>
          <CircularProgress size={50} />
          <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
            Ищем подходящие квесты по вашему запросу...
          </Typography>
        </Box>
      )}

      {results && (
        <Box>
          <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
            Найдено квестов: {results.length}
          </Typography>
          {results.length > 0 && (
            <MasonryGrid>
              {results.map((item, idx) => {
                // Извлекаем quest из структуры {quest: {...}, similarity_score: ...}
                const quest = item.quest || item;
                return (
                  <QuestCard
                    key={quest.id || idx}
                    quest={quest}
                    onAction={purchaseQuest}
                    onDetails={showQuestDetails}
                    actionLabel="Купить"
                  />
                );
              })}
            </MasonryGrid>
          )}
        </Box>
      )}

      <Dialog
        open={detailsDialogOpen}
        onClose={() => setDetailsDialogOpen(false)}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: {
            borderRadius: 3,
          }
        }}
      >
        <DialogTitle 
          sx={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            pb: 2,
            background: 'linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%)',
            color: 'white',
            borderRadius: '12px 12px 0 0',
          }}
        >
          <Typography variant="h4" fontWeight={700} sx={{ flex: 1 }}>
            {questDetails?.title || 'Детали квеста'}
          </Typography>
          <IconButton 
            onClick={() => setDetailsDialogOpen(false)}
            sx={{ color: 'white' }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent sx={{ p: 0 }}>
          {questDetails && (
            <Box sx={{ p: 3 }}>
              <Stack spacing={3}>
                {/* Описание */}
                <Paper 
                  elevation={0}
                  sx={{ 
                    p: 3, 
                    bgcolor: '#f8f9fa',
                    borderRadius: 2,
                    border: '1px solid #e9ecef'
                  }}
                >
                  <Typography 
                    variant="body1" 
                    sx={{ 
                      lineHeight: 1.8,
                      fontSize: '1.05rem',
                      color: '#495057'
                    }}
                  >
                    {questDetails.description}
                  </Typography>
                </Paper>

                {/* Награды */}
                <Box>
                  <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                    Награды
                  </Typography>
                  <Stack direction="row" spacing={2} flexWrap="wrap">
                    <Chip 
                      icon={<AttachMoneyIcon />}
                      label={`Цена: ${questDetails.price || 0} монет`}
                      sx={{ 
                        bgcolor: '#fee2e2',
                        color: '#991b1b',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#991b1b'
                        }
                      }}
                    />
                    <Chip 
                      icon={<StarIcon />}
                      label={`Награда: ${questDetails.reward_coin || 0} монет`}
                      sx={{ 
                        bgcolor: '#fef3c7',
                        color: '#92400e',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#92400e'
                        }
                      }}
                    />
                    <Chip 
                      icon={<TrendingUpIcon />}
                      label={`${questDetails.reward_xp || 0} XP`}
                      sx={{ 
                        bgcolor: '#dbeafe',
                        color: '#1e40af',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#1e40af'
                        }
                      }}
                    />
                  </Stack>
                </Box>

                <Divider />

                {/* Задачи */}
                {questDetails.tasks && questDetails.tasks.length > 0 && (
                  <Box>
                    <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                      Задачи ({questDetails.tasks.length})
                    </Typography>
                    <Stack spacing={2}>
                      {questDetails.tasks.map((task, idx) => (
                        <Paper
                          key={idx}
                          elevation={0}
                          sx={{
                            p: 2.5,
                            border: '1px solid #e9ecef',
                            borderRadius: 2,
                            bgcolor: '#ffffff',
                            transition: 'all 0.2s',
                            '&:hover': {
                              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                              borderColor: '#4a5568',
                            }
                          }}
                        >
                          <Stack direction="row" spacing={2} alignItems="flex-start">
                            <Box
                              sx={{
                                minWidth: 32,
                                height: 32,
                                borderRadius: '50%',
                                bgcolor: '#2d3748',
                                color: 'white',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                fontWeight: 700,
                                fontSize: '0.9rem',
                                flexShrink: 0,
                              }}
                            >
                              {idx + 1}
                            </Box>
                            <Box sx={{ flex: 1 }}>
                              <Typography 
                                variant="subtitle1" 
                                fontWeight={600}
                                gutterBottom
                                sx={{ mb: 1 }}
                              >
                                {task.title}
                              </Typography>
                              <Typography 
                                variant="body2" 
                                color="text.secondary"
                                sx={{ lineHeight: 1.6 }}
                              >
                                {task.description}
                              </Typography>
                            </Box>
                          </Stack>
                        </Paper>
                      ))}
                    </Stack>
                  </Box>
                )}
              </Stack>
            </Box>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 2.5, bgcolor: '#f8f9fa', borderRadius: '0 0 12px 12px' }}>
          <Button 
            onClick={() => setDetailsDialogOpen(false)}
            sx={{ minWidth: 100 }}
          >
            Закрыть
          </Button>
          {questDetails && (
            <Button 
              variant="contained" 
              onClick={() => {
                purchaseQuest(questDetails);
                setDetailsDialogOpen(false);
              }}
              startIcon={<AttachMoneyIcon />}
              sx={{ 
                minWidth: 150,
                bgcolor: '#2d3748',
                '&:hover': {
                  bgcolor: '#1a202c',
                }
              }}
            >
              Купить квест
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default SearchQuestsTab;
