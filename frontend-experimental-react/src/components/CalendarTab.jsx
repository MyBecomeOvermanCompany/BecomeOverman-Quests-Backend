import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Card, 
  CircularProgress,
  Stack,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Chip,
} from '@mui/material';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import CloseIcon from '@mui/icons-material/Close';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

function CalendarTab() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [eventDialogOpen, setEventDialogOpen] = useState(false);

  useEffect(() => {
    loadCalendarData();
  }, []);

  const loadCalendarData = async () => {
    setLoading(true);
    try {
      const result = await apiCall('/quests/my-quests-with-details');

      if (result.success && result.data && Array.isArray(result.data)) {
        const calendarEvents = mapQuestsToCalendarEvents(result.data);
        setEvents(calendarEvents);
      } else {
        showToast('Ошибка загрузки данных календаря', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки данных календаря', 'error');
    } finally {
      setLoading(false);
    }
  };

  const mapQuestsToCalendarEvents = (quests) => {
    const events = [];
    
    quests.forEach(quest => {
      if (quest.tasks && Array.isArray(quest.tasks)) {
        quest.tasks.forEach(task => {
          if (task.due_date || task.scheduled_start) {
            const dueDate = task.due_date || task.scheduled_start;
            const endDate = task.scheduled_end || dueDate;
            const now = new Date();
            
            let color = '#000000';
            if (task.completed || task.status === 'completed') {
              color = '#10b981';
            } else if (new Date(dueDate) < now) {
              color = '#ef4444';
            } else {
              const daysUntilDue = Math.ceil((new Date(dueDate) - now) / (1000 * 60 * 60 * 24));
              if (daysUntilDue <= 1) {
                color = '#f59e0b';
              }
            }

            events.push({
              id: `task-${task.id}`,
              title: task.title,
              start: dueDate,
              end: endDate,
              backgroundColor: color,
              borderColor: color,
              textColor: '#ffffff',
              extendedProps: {
                questId: quest.id,
                taskId: task.id,
                questTitle: quest.title,
                taskTitle: task.title,
                taskDescription: task.description,
                completed: task.completed || task.status === 'completed',
                category: quest.category,
              },
            });
          }
        });
      }
    });

    return events;
  };

  const handleEventClick = (info) => {
    setSelectedEvent(info.event.extendedProps);
    setEventDialogOpen(true);
  };

  const completeTaskFromCalendar = async () => {
    if (!selectedEvent || !selectedEvent.questId || !selectedEvent.taskId) return;

    try {
      const result = await apiCall(`/quests/${selectedEvent.questId}/${selectedEvent.taskId}/complete`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Задача выполнена!', 'success');
        setEventDialogOpen(false);
        loadCalendarData();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка выполнения задачи', 'error');
    }
  };

  const generateAISchedule = async () => {
    setLoading(true);
    try {
      const result = await apiCall('/quests/schedule', {
        method: 'POST',
      });

      if (result.success) {
        showToast('AI-расписание сгенерировано!', 'success');
        loadCalendarData();
      } else {
        showToast('Ошибка генерации расписания: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
        setLoading(false);
      }
    } catch (error) {
      showToast('Ошибка генерации расписания', 'error');
      setLoading(false);
    }
  };

  return (
    <Box>
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 4 }}>
        <CalendarTodayIcon sx={{ fontSize: 48, color: 'primary.main' }} />
        <Box>
          <Typography variant="h3" gutterBottom sx={{ fontWeight: 700 }}>
            Мое расписание
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Управление задачами и квестами
          </Typography>
        </Box>
      </Stack>

      <Button 
        variant="contained" 
        onClick={generateAISchedule} 
        disabled={loading} 
        sx={{ mb: 3 }}
        startIcon={<AutoAwesomeIcon />}
      >
        AI-расписание
      </Button>

      {loading && events.length === 0 ? (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress size={60} />
        </Box>
      ) : (
        <Card sx={{ p: 3 }}>
          <FullCalendar
            plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin]}
            initialView="dayGridMonth"
            events={events}
            eventClick={handleEventClick}
            headerToolbar={{
              left: 'prev,next today',
              center: 'title',
              right: 'dayGridMonth,timeGridWeek,timeGridDay',
            }}
            height="auto"
            eventDisplay="block"
            eventTextColor="#ffffff"
            eventTimeFormat={{
              hour: '2-digit',
              minute: '2-digit',
            }}
            slotMinTime="00:00:00"
            slotMaxTime="24:00:00"
            allDaySlot={true}
            locale="ru"
          />
        </Card>
      )}

      <Dialog
        open={eventDialogOpen}
        onClose={() => setEventDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6" fontWeight={600}>
            {selectedEvent?.taskTitle || 'Задача'}
          </Typography>
          <IconButton onClick={() => setEventDialogOpen(false)}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          {selectedEvent && (
            <Stack spacing={2}>
              <Box>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Квест
                </Typography>
                <Typography variant="body1" fontWeight={600}>
                  {selectedEvent.questTitle}
                </Typography>
              </Box>
              {selectedEvent.taskDescription && (
                <Box>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Описание
                  </Typography>
                  <Typography variant="body1">
                    {selectedEvent.taskDescription}
                  </Typography>
                </Box>
              )}
              {selectedEvent.completed && (
                <Chip 
                  icon={<CheckCircleIcon />}
                  label="Выполнено" 
                  color="success"
                />
              )}
            </Stack>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEventDialogOpen(false)}>
            Закрыть
          </Button>
          {selectedEvent && !selectedEvent.completed && (
            <Button 
              variant="contained" 
              onClick={completeTaskFromCalendar}
            >
              Выполнить задачу
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default CalendarTab;
