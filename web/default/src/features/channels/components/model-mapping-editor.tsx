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
import { useState, useEffect, useRef } from 'react'
import { Code, Table, Plus, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  parseModelMappingRows,
  serializeModelMappingRows,
  shouldSyncExternalMappingValue,
  type ModelMappingRow,
} from './model-mapping-editor-utils'

type ModelMappingEditorProps = {
  value: string
  onChange: (value: string) => void
  disabled?: boolean
}

export function ModelMappingEditor({
  value,
  onChange,
  disabled = false,
}: ModelMappingEditorProps) {
  const { t } = useTranslation()
  const [mode, setMode] = useState<'visual' | 'json'>('visual')
  const [rows, setRows] = useState<ModelMappingRow[]>([])
  const [jsonValue, setJsonValue] = useState(value)
  const nextRowId = useRef(0)
  const lastEmittedValue = useRef<string | null>(null)

  const createRowId = () => {
    const id = `mapping-${nextRowId.current}`
    nextRowId.current += 1
    return id
  }

  const syncRowsFromValue = (nextValue: string) => {
    const nextRows = parseModelMappingRows(nextValue, createRowId)
    if (nextRows !== null) setRows(nextRows)
  }

  const emitChange = (nextValue: string) => {
    setJsonValue(nextValue)
    lastEmittedValue.current = nextValue
    onChange(nextValue)
  }

  // Parse JSON to rows when value changes externally
  useEffect(() => {
    setJsonValue(value)
    if (shouldSyncExternalMappingValue(value, lastEmittedValue.current)) {
      syncRowsFromValue(value)
    }
    lastEmittedValue.current = null
    // The only external input is `value`; callbacks use stable ref-backed IDs.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value])

  const handleAddRow = () => {
    const newRow: ModelMappingRow = {
      id: createRowId(),
      from: '',
      to: '',
    }
    const updatedRows = [...rows, newRow]
    setRows(updatedRows)
  }

  const handleDeleteRow = (id: string) => {
    const updatedRows = rows.filter((row) => row.id !== id)
    setRows(updatedRows)
    emitChange(serializeModelMappingRows(updatedRows))
  }

  const handleRowChange = (
    id: string,
    field: 'from' | 'to',
    newValue: string
  ) => {
    const updatedRows = rows.map((row) =>
      row.id === id ? { ...row, [field]: newValue } : row
    )
    setRows(updatedRows)
    emitChange(serializeModelMappingRows(updatedRows))
  }

  const handleJsonChange = (newJson: string) => {
    emitChange(newJson)
    syncRowsFromValue(newJson)
  }

  const handleFillTemplate = () => {
    const template = JSON.stringify(
      { 'gpt-3.5-turbo': 'gpt-3.5-turbo-0125' },
      null,
      2
    )
    emitChange(template)
    syncRowsFromValue(template)
  }

  const toggleMode = () => {
    if (mode === 'visual') {
      // Switching to JSON mode: sync rows to JSON
      const json = serializeModelMappingRows(rows)
      emitChange(json)
      setMode('json')
    } else {
      // Switching to visual mode: sync JSON to rows
      syncRowsFromValue(jsonValue)
      setMode('visual')
    }
  }

  return (
    <div className='space-y-2'>
      <div className='flex items-center justify-between'>
        <div className='flex gap-2'>
          <Button
            type='button'
            variant='outline'
            size='sm'
            onClick={toggleMode}
            disabled={disabled}
          >
            {mode === 'visual' ? (
              <>
                <Code className='mr-2 h-4 w-4' />
                {t('JSON Mode')}
              </>
            ) : (
              <>
                <Table className='mr-2 h-4 w-4' />
                {t('Visual Mode')}
              </>
            )}
          </Button>
          <Button
            type='button'
            variant='link'
            size='sm'
            className='h-auto p-0'
            onClick={handleFillTemplate}
            disabled={disabled}
          >
            {t('Fill Template')}
          </Button>
        </div>
      </div>

      {mode === 'visual' ? (
        <div className='space-y-2'>
          {rows.length > 0 ? (
            <div className='space-y-2'>
              <div className='grid grid-cols-[1fr_1fr_auto] gap-2 text-sm font-medium'>
                <div>{t('Original Model')}</div>
                <div>{t('Replacement Model')}</div>
                <div className='w-10'></div>
              </div>
              {rows.map((row) => (
                <div
                  key={row.id}
                  className='grid grid-cols-[1fr_1fr_auto] gap-2'
                >
                  <Input
                    value={row.from}
                    onChange={(e) =>
                      handleRowChange(row.id, 'from', e.target.value)
                    }
                    placeholder='gpt-3.5-turbo'
                    disabled={disabled}
                  />
                  <Input
                    value={row.to}
                    onChange={(e) =>
                      handleRowChange(row.id, 'to', e.target.value)
                    }
                    placeholder='gpt-3.5-turbo-0125'
                    disabled={disabled}
                  />
                  <Button
                    type='button'
                    variant='ghost'
                    size='icon'
                    onClick={() => handleDeleteRow(row.id)}
                    disabled={disabled}
                    className='h-10 w-10'
                  >
                    <Trash2 className='h-4 w-4' />
                  </Button>
                </div>
              ))}
            </div>
          ) : (
            <div className='text-muted-foreground flex h-24 items-center justify-center rounded-md border border-dashed text-sm'>
              {t(
                'No model mappings configured. Click "Add Mapping" to get started.'
              )}
            </div>
          )}
          <Button
            type='button'
            variant='outline'
            size='sm'
            onClick={handleAddRow}
            disabled={disabled}
            className='w-full'
          >
            <Plus className='mr-2 h-4 w-4' />
            {t('Add Mapping')}
          </Button>
        </div>
      ) : (
        <Textarea
          value={jsonValue}
          onChange={(e) => handleJsonChange(e.target.value)}
          placeholder={t('{"original-model": "replacement-model"}')}
          disabled={disabled}
          rows={8}
          className={cn('font-mono text-sm')}
        />
      )}
    </div>
  )
}
