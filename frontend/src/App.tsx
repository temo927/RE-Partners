import PackSizes from './components/PackSizes'
import Calculator from './components/Calculator'
import './App.css'

function App() {
  return (
    <div className="app">
      <h1 className="title">Order Packs Calculator</h1>
      <PackSizes />
      <Calculator />
    </div>
  )
}

export default App
