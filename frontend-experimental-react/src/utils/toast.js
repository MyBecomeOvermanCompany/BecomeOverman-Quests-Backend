import { toast } from 'react-toastify';

export const showToast = (message, type = 'info') => {
  const options = {
    position: 'top-right',
    autoClose: 3000,
    hideProgressBar: false,
    closeOnClick: true,
    pauseOnHover: true,
    draggable: true,
    progress: undefined,
    style: {
      borderRadius: '12px',
      fontFamily: "'Inter', sans-serif",
    },
  };

  switch (type) {
    case 'success':
      toast.success(message, {
        ...options,
        icon: '✅',
      });
      break;
    case 'error':
      toast.error(message, {
        ...options,
        icon: '❌',
      });
      break;
    case 'warning':
      toast.warning(message, {
        ...options,
        icon: '⚠️',
      });
      break;
    case 'info':
    default:
      toast.info(message, {
        ...options,
        icon: 'ℹ️',
      });
      break;
  }
};

export default showToast;
