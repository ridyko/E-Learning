import Swal from 'sweetalert2'

export const showSuccess = (title, text = '') => {
  return Swal.fire({
    title: title,
    text: text,
    icon: 'success',
    confirmButtonColor: '#1d4ed8',
    confirmButtonText: 'OK, Baik',
    background: '#ffffff',
    color: '#0f172a',
    customClass: {
      popup: 'swal-custom-popup',
      title: 'swal-custom-title'
    }
  })
}

export const showError = (title, text = '') => {
  return Swal.fire({
    title: title,
    text: text,
    icon: 'error',
    confirmButtonColor: '#dc2626',
    confirmButtonText: 'Tutup',
    background: '#ffffff',
    color: '#0f172a'
  })
}

export const showWarning = (title, text = '') => {
  return Swal.fire({
    title: title,
    text: text,
    icon: 'warning',
    confirmButtonColor: '#d97706',
    confirmButtonText: 'Mengerti',
    background: '#ffffff',
    color: '#0f172a'
  })
}

export const showConfirm = async (title, text = '', confirmText = 'Ya, Lanjutkan', callback = null) => {
  let confirmBtnLabel = 'Ya, Lanjutkan'
  let cb = null

  if (typeof confirmText === 'function') {
    cb = confirmText
  } else {
    confirmBtnLabel = confirmText
    if (typeof callback === 'function') cb = callback
  }

  const result = await Swal.fire({
    title: title,
    text: text,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    cancelButtonColor: '#64748b',
    confirmButtonText: confirmBtnLabel,
    cancelButtonText: 'Batal',
    background: '#ffffff',
    color: '#0f172a'
  })

  if (result.isConfirmed && cb) {
    cb()
  }

  return result.isConfirmed
}

export const showToast = (title, icon = 'success') => {
  const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 3000,
    timerProgressBar: true,
    didOpen: (toast) => {
      toast.onmouseenter = Swal.stopTimer
      toast.onmouseleave = Swal.resumeTimer
    }
  })
  return Toast.fire({
    icon: icon,
    title: title
  })
}

export default Swal
