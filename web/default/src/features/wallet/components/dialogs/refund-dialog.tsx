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
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { requestOpenPaymentRefund, isApiSuccess } from '../../api'
import type { TopupRecord } from '../../types'

interface RefundDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  record: TopupRecord | null
  onSuccess?: () => void
}

export function RefundDialog({
  open,
  onOpenChange,
  record,
  onSuccess,
}: RefundDialogProps) {
  const { t } = useTranslation()
  const [refundAmount, setRefundAmount] = useState('')
  const [reason, setReason] = useState('')
  const [processing, setProcessing] = useState(false)

  if (!record) return null

  const maxRefundable = record.money - (record.refund_amount || 0)

  const handleOpen = (isOpen: boolean) => {
    if (isOpen) {
      setRefundAmount(maxRefundable.toFixed(2))
      setReason('')
    }
    onOpenChange(isOpen)
  }

  const handleSubmit = async () => {
    const amount = parseFloat(refundAmount)
    if (isNaN(amount) || amount <= 0) {
      toast.error(t('Please enter a valid refund amount'))
      return
    }
    if (amount > maxRefundable) {
      toast.error(
        t('Refund amount cannot exceed {{max}}', {
          max: maxRefundable.toFixed(2),
        })
      )
      return
    }

    setProcessing(true)
    try {
      const res = await requestOpenPaymentRefund({
        trade_no: record.trade_no,
        refund_amount: amount,
        reason: reason || undefined,
      })

      if (isApiSuccess(res)) {
        toast.success(t('Refund initiated successfully'))
        onOpenChange(false)
        onSuccess?.()
      } else {
        toast.error(
          (res as { message?: string }).message || t('Refund request failed')
        )
      }
    } catch {
      toast.error(t('Refund request failed'))
    } finally {
      setProcessing(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t('Initiate Refund')}</DialogTitle>
          <DialogDescription>
            {t('Refund will deduct equivalent quota from the user')}
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-2">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <Label className="text-muted-foreground text-xs">
                {t('Order No')}
              </Label>
              <div className="font-mono text-xs">{record.trade_no}</div>
            </div>
            <div>
              <Label className="text-muted-foreground text-xs">
                {t('Payment Amount')}
              </Label>
              <div className="font-medium">¥{record.money.toFixed(2)}</div>
            </div>
            <div>
              <Label className="text-muted-foreground text-xs">
                {t('Already Refunded')}
              </Label>
              <div className="font-medium">
                ¥{(record.refund_amount || 0).toFixed(2)}
              </div>
            </div>
            <div>
              <Label className="text-muted-foreground text-xs">
                {t('Max Refundable')}
              </Label>
              <div className="font-medium text-green-600">
                ¥{maxRefundable.toFixed(2)}
              </div>
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="refund-amount">{t('Refund Amount (¥)')}</Label>
            <Input
              id="refund-amount"
              type="number"
              step="0.01"
              min="0.01"
              max={maxRefundable}
              value={refundAmount}
              onChange={(e) => setRefundAmount(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="refund-reason">
              {t('Reason')} ({t('Optional')})
            </Label>
            <Textarea
              id="refund-reason"
              rows={2}
              placeholder={t('Admin refund')}
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={processing}
          >
            {t('Cancel')}
          </Button>
          <Button onClick={handleSubmit} disabled={processing}>
            {processing ? t('Processing...') : t('Confirm Refund')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
