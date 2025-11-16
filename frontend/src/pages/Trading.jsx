import React, { useState } from 'react'
import { useMutation, useQueryClient } from 'react-query'
import { tradeService } from '../services/api'
import { useAuth } from '../contexts/AuthContext'

const ASSETS = [
  { symbol: 'BTC', name: 'Bitcoin', price: 45000, type: 'SPOT' },
  { symbol: 'ETH', name: 'Ethereum', price: 2500, type: 'SPOT' },
  { symbol: 'SOL', name: 'Solana', price: 100, type: 'SPOT' },
  { symbol: 'USDC', name: 'USD Coin', price: 1.0, type: 'SPOT' },
  { symbol: 'USDT', name: 'Tether', price: 1.0, type: 'SPOT' },
  { symbol: 'BTC-PERP', name: 'Bitcoin Perpetual', price: 45000, type: 'FUTURES' },
  { symbol: 'ETH-PERP', name: 'Ethereum Perpetual', price: 2500, type: 'FUTURES' },
  { symbol: 'SOL-PERP', name: 'Solana Perpetual', price: 100, type: 'FUTURES' },
]

const Trading = () => {
  const { user } = useAuth()
  const queryClient = useQueryClient()
  const [selectedAsset, setSelectedAsset] = useState(ASSETS[0])
  const [tradeType, setTradeType] = useState('BUY')
  const [quantity, setQuantity] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  const buyMutation = useMutation(
    ({ asset_symbol, quantity, price }) => tradeService.buy(asset_symbol, quantity, price),
    {
      onSuccess: (data) => {
        setMessage(`Successfully bought ${quantity} ${selectedAsset.symbol}`)
        setQuantity('')
        setError('')
        queryClient.invalidateQueries('portfolio')
        setTimeout(() => window.location.reload(), 1000)
      },
      onError: (err) => {
        setError(err.error || 'Failed to execute trade')
        setMessage('')
      },
    }
  )

  const sellMutation = useMutation(
    ({ asset_symbol, quantity, price }) => tradeService.sell(asset_symbol, quantity, price),
    {
      onSuccess: (data) => {
        setMessage(`Successfully sold ${quantity} ${selectedAsset.symbol}`)
        setQuantity('')
        setError('')
        queryClient.invalidateQueries('portfolio')
        setTimeout(() => window.location.reload(), 1000)
      },
      onError: (err) => {
        setError(err.error || 'Failed to execute trade')
        setMessage('')
      },
    }
  )

  const handleTrade = (e) => {
    e.preventDefault()
    setMessage('')
    setError('')

    const qty = parseFloat(quantity)
    if (isNaN(qty) || qty <= 0) {
      setError('Please enter a valid quantity')
      return
    }

    const tradeData = {
      asset_symbol: selectedAsset.symbol,
      quantity: qty,
      price: selectedAsset.price,
    }

    if (tradeType === 'BUY') {
      buyMutation.mutate(tradeData)
    } else {
      sellMutation.mutate(tradeData)
    }
  }

  const totalCost = quantity ? (parseFloat(quantity) * selectedAsset.price).toFixed(2) : '0.00'

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Trading</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
            <h2 className="text-xl font-semibold mb-6">Place Order</h2>

            {message && (
              <div className="mb-4 bg-green-500 bg-opacity-10 border border-green-500 text-green-500 px-4 py-3 rounded">
                {message}
              </div>
            )}

            {error && (
              <div className="mb-4 bg-red-500 bg-opacity-10 border border-red-500 text-red-500 px-4 py-3 rounded">
                {error}
              </div>
            )}

            <form onSubmit={handleTrade} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">Asset</label>
                <select
                  value={selectedAsset.symbol}
                  onChange={(e) =>
                    setSelectedAsset(ASSETS.find((a) => a.symbol === e.target.value))
                  }
                  className="w-full px-3 py-2 border border-gray-600 rounded-lg bg-slate-700 text-white focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                >
                  {ASSETS.map((asset) => (
                    <option key={asset.symbol} value={asset.symbol}>
                      {asset.name} ({asset.symbol}) - ${asset.price.toFixed(2)}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex space-x-4">
                <button
                  type="button"
                  onClick={() => setTradeType('BUY')}
                  className={`flex-1 py-3 rounded-lg font-medium ${
                    tradeType === 'BUY'
                      ? 'bg-green-600 text-white'
                      : 'bg-slate-700 text-gray-300 hover:bg-slate-600'
                  }`}
                >
                  Buy
                </button>
                <button
                  type="button"
                  onClick={() => setTradeType('SELL')}
                  className={`flex-1 py-3 rounded-lg font-medium ${
                    tradeType === 'SELL'
                      ? 'bg-red-600 text-white'
                      : 'bg-slate-700 text-gray-300 hover:bg-slate-600'
                  }`}
                >
                  Sell
                </button>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">Price</label>
                <input
                  type="text"
                  value={`$${selectedAsset.price.toFixed(2)}`}
                  disabled
                  className="w-full px-3 py-2 border border-gray-600 rounded-lg bg-slate-700 text-gray-400"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">Quantity</label>
                <input
                  type="number"
                  step="0.00001"
                  value={quantity}
                  onChange={(e) => setQuantity(e.target.value)}
                  placeholder="0.00"
                  className="w-full px-3 py-2 border border-gray-600 rounded-lg bg-slate-700 text-white focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div className="bg-slate-700 rounded-lg p-4">
                <div className="flex justify-between mb-2">
                  <span className="text-gray-400">Total Cost:</span>
                  <span className="font-semibold">${totalCost}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Available Balance:</span>
                  <span className="font-semibold">${user?.balance?.toFixed(2)}</span>
                </div>
              </div>

              <button
                type="submit"
                disabled={buyMutation.isLoading || sellMutation.isLoading}
                className={`w-full py-3 rounded-lg font-medium ${
                  tradeType === 'BUY'
                    ? 'bg-green-600 hover:bg-green-700'
                    : 'bg-red-600 hover:bg-red-700'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {buyMutation.isLoading || sellMutation.isLoading
                  ? 'Processing...'
                  : `${tradeType} ${selectedAsset.symbol}`}
              </button>
            </form>
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
            <h3 className="text-lg font-semibold mb-4">Market Overview</h3>
            <div className="space-y-3">
              {ASSETS.filter((a) => a.type === 'SPOT').map((asset) => (
                <div
                  key={asset.symbol}
                  className="flex justify-between items-center p-3 bg-slate-700 rounded-lg cursor-pointer hover:bg-slate-600"
                  onClick={() => setSelectedAsset(asset)}
                >
                  <div>
                    <div className="font-medium">{asset.symbol}</div>
                    <div className="text-sm text-gray-400">{asset.name}</div>
                  </div>
                  <div className="text-right">
                    <div className="font-semibold">${asset.price.toFixed(2)}</div>
                    <div className="text-sm text-green-400">+0.5%</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
            <h3 className="text-lg font-semibold mb-4">Futures</h3>
            <div className="space-y-3">
              {ASSETS.filter((a) => a.type === 'FUTURES').map((asset) => (
                <div
                  key={asset.symbol}
                  className="flex justify-between items-center p-3 bg-slate-700 rounded-lg cursor-pointer hover:bg-slate-600"
                  onClick={() => setSelectedAsset(asset)}
                >
                  <div>
                    <div className="font-medium">{asset.symbol}</div>
                    <div className="text-sm text-gray-400">Perpetual</div>
                  </div>
                  <div className="text-right">
                    <div className="font-semibold">${asset.price.toFixed(2)}</div>
                    <div className="text-sm text-green-400">+0.5%</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Trading
