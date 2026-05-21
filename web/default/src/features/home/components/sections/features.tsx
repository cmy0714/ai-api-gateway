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
import {
  Zap,
  Shield,
  Globe,
  BarChart3,
  Gauge,
  DollarSign,
  Users,
  Clock,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { AnimateInView } from '@/components/animate-in-view'

interface FeaturesProps {
  className?: string
}

export function Features(_props: FeaturesProps) {
  const { t } = useTranslation()

  const features = [
    {
      id: 'fast',
      num: '01',
      title: t('Blazing Performance'),
      desc: t(
        'Intelligent routing and edge caching deliver sub-200ms latency across all models'
      ),
      span: 'md:col-span-2',
      icon: <Zap className='size-4 text-orange-400' />,
      visual: (
        <div className='mt-4 grid grid-cols-3 gap-2'>
          {[
            { name: 'GPT-4o', ms: '142ms' },
            { name: 'Claude', ms: '156ms' },
            { name: 'Gemini', ms: '93ms' },
            { name: 'DeepSeek', ms: '118ms' },
            { name: 'Qwen', ms: '105ms' },
            { name: 'Llama', ms: '88ms' },
          ].map((model) => (
            <div
              key={model.name}
              className='border-border/30 bg-muted/20 flex items-center justify-between rounded-lg border px-3 py-2 transition-colors duration-300 hover:border-orange-500/30 hover:bg-orange-500/5'
            >
              <span className='text-muted-foreground text-xs'>
                {model.name}
              </span>
              <span className='font-mono text-[10px] text-emerald-600 dark:text-emerald-400'>
                {model.ms}
              </span>
            </div>
          ))}
        </div>
      ),
    },
    {
      id: 'secure',
      num: '02',
      title: t('Enterprise Security'),
      desc: t(
        'End-to-end encryption, role-based access control, and full audit trail'
      ),
      span: 'md:col-span-1',
      icon: <Shield className='size-4 text-emerald-400' />,
      visual: (
        <div className='mt-4 flex items-center justify-center'>
          <div className='relative'>
            <div className='flex size-16 items-center justify-center rounded-2xl border border-emerald-500/20 bg-emerald-500/5'>
              <Shield
                className='size-7 text-emerald-500/70'
                strokeWidth={1.5}
              />
            </div>
            <div className='absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full bg-emerald-500'>
              <svg
                className='size-2.5 text-white'
                fill='none'
                viewBox='0 0 24 24'
                stroke='currentColor'
                strokeWidth={3}
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  d='m4.5 12.75 6 6 9-13.5'
                />
              </svg>
            </div>
          </div>
        </div>
      ),
    },
    {
      id: 'global',
      num: '03',
      title: t('Smart Load Balancing'),
      desc: t(
        'Automatic failover and intelligent traffic distribution across providers'
      ),
      span: 'md:col-span-1',
      icon: <Globe className='size-4 text-violet-400' />,
      visual: (
        <div className='mt-4 space-y-2'>
          {[
            t('Auto Failover'),
            t('Traffic Routing'),
            t('Health Monitoring'),
          ].map((step, i) => (
            <div key={step} className='flex items-center gap-2'>
              <div
                className={`flex size-6 items-center justify-center rounded-full text-[10px] font-bold ${
                  i === 1
                    ? 'border border-orange-500/30 bg-orange-500/20 text-orange-500'
                    : 'border-border/40 bg-muted text-muted-foreground border'
                }`}
              >
                {i + 1}
              </div>
              <div className='bg-border/40 h-px flex-1' />
              <span className='text-muted-foreground text-xs'>{step}</span>
            </div>
          ))}
        </div>
      ),
    },
    {
      id: 'analytics',
      num: '04',
      title: t('Real-time Analytics'),
      desc: t(
        'Comprehensive dashboards for usage tracking, cost analysis, and performance insights'
      ),
      span: 'md:col-span-2',
      icon: <BarChart3 className='size-4 text-amber-400' />,
      visual: (
        <div className='mt-4 flex items-center gap-4'>
          <div className='flex flex-1 items-end gap-1'>
            {[40, 65, 45, 80, 55, 90, 70, 85].map((h, i) => (
              <div
                key={i}
                className='flex-1 rounded-sm bg-gradient-to-t from-orange-500/40 to-orange-500/10 transition-all duration-300 hover:from-orange-500/60 hover:to-orange-500/20'
                style={{ height: `${h}%`, minHeight: `${h * 0.4}px` }}
              />
            ))}
          </div>
          <div className='text-muted-foreground space-y-1 text-[10px]'>
            <div className='flex items-center gap-1.5'>
              <span className='inline-block size-2 rounded-full bg-orange-500/60' />
              {t('Requests')}
            </div>
            <div className='flex items-center gap-1.5'>
              <span className='inline-block size-2 rounded-full bg-emerald-500/60' />
              {t('Tokens')}
            </div>
          </div>
        </div>
      ),
    },
  ]

  const additionalFeatures = [
    {
      icon: <Gauge className='size-5' strokeWidth={1.5} />,
      title: t('Auto Scaling'),
      desc: t('Handles millions of requests with zero configuration'),
    },
    {
      icon: <DollarSign className='size-5' strokeWidth={1.5} />,
      title: t('Transparent Billing'),
      desc: t('Pay only for what you use with real-time cost tracking'),
    },
    {
      icon: <Users className='size-5' strokeWidth={1.5} />,
      title: t('Team Management'),
      desc: t('Granular permissions and multi-tenant workspace support'),
    },
    {
      icon: <Clock className='size-5' strokeWidth={1.5} />,
      title: t('24/7 Reliability'),
      desc: t('99.9% uptime SLA backed by global infrastructure'),
    },
  ]

  return (
    <section className='relative z-10 px-6 py-24 md:py-32'>
      <div className='mx-auto max-w-6xl'>
        <AnimateInView className='mb-16 max-w-lg'>
          <p className='text-muted-foreground mb-3 text-xs font-medium tracking-widest uppercase'>
            {t('Why Choose Us')}
          </p>
          <h2 className='text-2xl leading-tight font-bold tracking-tight md:text-3xl'>
            {t('Built for performance,')}
            <br />
            {t('designed for growth')}
          </h2>
        </AnimateInView>

        {/* Bento grid */}
        <div className='border-border/40 bg-border/40 grid gap-px overflow-hidden rounded-2xl border md:grid-cols-3'>
          {features.map((f, i) => (
            <AnimateInView
              key={f.id}
              delay={i * 100}
              animation='scale-in'
              className={`bg-background group hover:bg-muted/20 p-7 transition-colors duration-300 md:p-8 ${f.span}`}
            >
              <div className='mb-3 flex items-center gap-3'>
                <span className='border-border/40 bg-muted text-muted-foreground flex size-7 items-center justify-center rounded-md border text-[10px] font-semibold tabular-nums'>
                  {f.num}
                </span>
                <h3 className='text-sm font-semibold'>{f.title}</h3>
              </div>
              <p className='text-muted-foreground text-sm leading-relaxed'>
                {f.desc}
              </p>
              {f.visual}
            </AnimateInView>
          ))}
        </div>

        {/* Additional features row */}
        <div className='mt-12 grid grid-cols-2 gap-8 md:grid-cols-4 md:gap-12'>
          {additionalFeatures.map((f, i) => (
            <AnimateInView
              key={f.title}
              delay={i * 100}
              animation='fade-up'
              className='flex flex-col items-center text-center'
            >
              <div className='text-muted-foreground border-border/50 bg-muted/30 group-hover:text-foreground mb-3 flex size-12 items-center justify-center rounded-xl border transition-colors'>
                {f.icon}
              </div>
              <h3 className='mb-1.5 text-sm font-semibold'>{f.title}</h3>
              <p className='text-muted-foreground max-w-[200px] text-xs leading-relaxed'>
                {f.desc}
              </p>
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}
