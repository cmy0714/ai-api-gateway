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
import { Link } from '@tanstack/react-router'
import { ArrowRight, Sparkles } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { HeroTerminalDemo } from '../hero-terminal-demo'

interface HeroProps {
  className?: string
  isAuthenticated?: boolean
}

export function Hero(props: HeroProps) {
  const { t } = useTranslation()

  return (
    <section className='relative z-10 flex flex-col items-center overflow-hidden px-6 pt-32 pb-20 md:pt-40 md:pb-32'>
      {/* Warm radial gradient background */}
      <div
        aria-hidden
        className='pointer-events-none absolute inset-0 -z-10 opacity-30 dark:opacity-[0.15]'
        style={{
          background: [
            'radial-gradient(ellipse 70% 50% at 50% 0%, oklch(0.75 0.20 50 / 90%) 0%, transparent 70%)',
            'radial-gradient(ellipse 40% 40% at 0% 50%, oklch(0.80 0.15 60 / 50%) 0%, transparent 70%)',
            'radial-gradient(ellipse 40% 40% at 100% 50%, oklch(0.70 0.18 40 / 40%) 0%, transparent 70%)',
          ].join(', '),
        }}
      />
      {/* Subtle dot pattern */}
      <div
        aria-hidden
        className='absolute inset-0 -z-10 bg-[radial-gradient(circle,var(--border)_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_30%,black_20%,transparent_100%)] bg-[size:2rem_2rem] opacity-[0.06]'
      />

      {/* Floating decorative orbs */}
      <div
        aria-hidden
        className='landing-animate-fade-in pointer-events-none absolute top-20 -left-20 size-72 rounded-full bg-orange-400/10 blur-3xl'
      />
      <div
        aria-hidden
        className='landing-animate-fade-in pointer-events-none absolute -right-16 bottom-10 size-60 rounded-full bg-amber-400/10 blur-3xl'
      />

      <div className='flex max-w-4xl flex-col items-center text-center'>
        {/* Badge */}
        <div
          className='landing-animate-fade-up mb-8 inline-flex items-center gap-2 rounded-full border border-orange-500/20 bg-orange-500/5 px-4 py-1.5 text-xs font-medium text-orange-600 dark:text-orange-400'
          style={{ animationDelay: '0ms' }}
        >
          <Sparkles className='size-3.5' />
          {t('Next-Gen AI Platform')}
        </div>

        <h1
          className='landing-animate-fade-up text-[clamp(2.2rem,6vw,4rem)] leading-[1.1] font-bold tracking-tight'
          style={{ animationDelay: '60ms' }}
        >
          {t('One Platform')}
          <br />
          <span className='bg-gradient-to-r from-orange-400 via-amber-500 to-orange-600 bg-clip-text text-transparent'>
            {t('Infinite AI Possibilities')}
          </span>
        </h1>
        <p
          className='landing-animate-fade-up text-muted-foreground/80 mt-6 max-w-xl text-base leading-relaxed opacity-0 md:text-lg'
          style={{ animationDelay: '120ms' }}
        >
          {t(
            'Access 50+ top AI models through a single API. Enterprise-grade performance, intelligent routing, and transparent billing — everything you need to power world-class AI applications.'
          )}
        </p>
        <div
          className='landing-animate-fade-up mt-10 flex items-center gap-4 opacity-0'
          style={{ animationDelay: '200ms' }}
        >
          {props.isAuthenticated ? (
            <Button
              size='lg'
              className='group rounded-xl bg-gradient-to-r from-orange-500 to-amber-500 px-8 text-white shadow-lg shadow-orange-500/25 transition-all hover:shadow-orange-500/40'
              render={<Link to='/dashboard' />}
            >
              {t('Go to Dashboard')}
              <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
            </Button>
          ) : (
            <>
              <Button
                size='lg'
                className='group rounded-xl bg-gradient-to-r from-orange-500 to-amber-500 px-8 text-white shadow-lg shadow-orange-500/25 transition-all hover:shadow-orange-500/40'
                render={<Link to='/sign-up' />}
              >
                {t('Start for Free')}
                <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
              </Button>
              <Button
                size='lg'
                variant='outline'
                className='border-border/50 hover:border-border hover:bg-muted/50 rounded-xl px-8'
                render={<Link to='/pricing' />}
              >
                {t('View Pricing')}
              </Button>
            </>
          )}
        </div>

        {/* Trust indicators */}
        <div
          className='landing-animate-fade-up mt-12 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 opacity-0'
          style={{ animationDelay: '320ms' }}
        >
          {[
            t('Enterprise SLA'),
            t('99.9% Uptime'),
            t('Global CDN'),
            t('SOC 2 Compliant'),
          ].map((item) => (
            <span
              key={item}
              className='text-muted-foreground/50 flex items-center gap-1.5 text-xs'
            >
              <span className='inline-block size-1 rounded-full bg-emerald-500' />
              {item}
            </span>
          ))}
        </div>
      </div>

      <div
        className='landing-animate-fade-up w-full opacity-0'
        style={{ animationDelay: '400ms' }}
      >
        <HeroTerminalDemo />
      </div>
    </section>
  )
}
