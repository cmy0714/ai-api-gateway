/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useState, useCallback, useRef } from 'react'
import i18next from 'i18next'
import { toast } from 'sonner'
import { requestOpenPayment, getOpenPaymentStatus, isApiSuccess } from '../api'
import type { OpenPaymentPayResponse } from '../types'

export function useOpenPayment() {
  const [processing, setProcessing] = useState(false)
  const [polling, setPolling] = useState(false)
  const [payResponse, setPayResponse] = useState<OpenPaymentPayResponse | null>(
    null
  )
  const [showQrDialog, setShowQrDialog] = useState(false)
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null)

  const stopPolling = useCallback(() => {
    if (pollingRef.current) {
      clearInterval(pollingRef.current)
      pollingRef.current = null
    }
    setPolling(false)
  }, [])

  const startPolling = useCallback(
    (tradeNo: string, onSuccess: () => void) => {
      setPolling(true)
      pollingRef.current = setInterval(async () => {
        try {
          const res = await getOpenPaymentStatus(tradeNo)
          if (isApiSuccess(res) && res.data?.status === 'success') {
            stopPolling()
            setShowQrDialog(false)
            setPayResponse(null)
            toast.success(i18next.t('Payment successful'))
            onSuccess()
          }
        } catch {
          // ignore polling errors
        }
      }, 2000)
    },
    [stopPolling]
  )

  const processOpenPayment = useCallback(
    async (
      amount: number,
      paymentMethod: string,
      onSuccess: () => void
    ): Promise<boolean> => {
      try {
        setProcessing(true)
        const response = await requestOpenPayment({
          amount,
          payment_method: paymentMethod,
        })

        if (!isApiSuccess(response) || !response.data) {
          toast.error(
            (response as { message?: string }).message ||
              i18next.t('Payment request failed')
          )
          return false
        }

        const data = response.data
        setPayResponse(data)

        if (data.pay_mode === 'qr' && data.qr_code_url) {
          setShowQrDialog(true)
          startPolling(data.trade_no, onSuccess)
          return true
        }

        if (
          (data.pay_mode === 'redirect' || data.pay_mode === 'params') &&
          data.pay_url
        ) {
          window.open(data.pay_url, '_blank')
          startPolling(data.trade_no, onSuccess)
          toast.success(i18next.t('Redirecting to payment page...'))
          return true
        }

        toast.error(i18next.t('Payment request failed'))
        return false
      } catch {
        toast.error(i18next.t('Payment request failed'))
        return false
      } finally {
        setProcessing(false)
      }
    },
    [startPolling]
  )

  const closeQrDialog = useCallback(() => {
    stopPolling()
    setShowQrDialog(false)
    setPayResponse(null)
  }, [stopPolling])

  return {
    processing,
    polling,
    payResponse,
    showQrDialog,
    processOpenPayment,
    closeQrDialog,
    stopPolling,
  }
}
