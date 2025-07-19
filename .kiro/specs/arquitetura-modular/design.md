# Design Document - Arquitetura Modular do Sistema

## Overview

Este documento define o design arquitetural do sistema de gestão empresarial, estabelecendo padrões, estruturas e diretrizes que devem ser seguidas em todas as implementações. O design é baseado na arquitetura bem-sucedida implementada no módulo de produtos e serve como referência para manter consistência e qualidade em todo o projeto.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  Pages (App Router)  │  Components (UI/Shared/Domain)       │
│  - /produtos         │  - Button, Input, Modal             │
│  - /clientes         │  - DataTable, Pagination            │
│  - /vendas           │  - ProdutoCard, ClienteForm          │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                     BUSINESS LOGIC LAYER                    │
├─────────────────────────────────────────────────────────────┤
│  Custom Hooks        │  Services                           │
│  - useProdutos       │  - produtosService                  │
│  - useClientes       │  - clientesService                  │
│  - useVendas         │  - vendasService                    │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      DATA ACCESS LAYER                      │
├─────────────────────────────────────────────────────────────┤
│  HTTP Client         │  Types & Validation                 │
│  - apiClient         │  - entities.ts                      │
│  - Error Handling    │  - api.ts                           │
│  - Request/Response  │  - common.ts                        │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                        EXTERNAL APIs                        │
├─────────────────────────────────────────────────────────────┤
│  REST APIs           │  Database                           │
│  - /api/produtos     │  - PostgreSQL                      │
│  - /api/clientes     │  - MongoDB                         │
│  - /api/vendas       │  - External APIs                   │
└─────────────────────────────────────────────────────────────┘
```

### Module Structure Pattern

Cada módulo deve seguir esta estrutura padrão:

```
/[modulo]
├── page.tsx                    # Página principal do módulo
├── [modulo].module.css         # CSS específico do módulo (opcional)
├── loading.tsx                 # Loading state (opcional)
├── error.tsx                   # Error boundary (opcional)
└── not-found.tsx              # 404 page (opcional)

/components/[modulo]
├── [Modulo]Card.tsx           # Card component
├── [Modulo]Table.tsx          # Table component
├── [Modulo]Form.tsx           # Form component
├── [Modulo]Modal.tsx          # Modal component
├── [Modulo]Filters.tsx        # Filters component
└── [Modulo]Stats.tsx          # Statistics component

/hooks
└── use[Modulo].ts             # Custom hook para o módulo

/lib/services
└── [modulo].service.ts        # Service layer

/types
└── entities.ts                # Adicionar interface do módulo
```

## Components and Interfaces

### 1. Page Component Pattern

```typescript
// Padrão para páginas de módulo
export default function [Modulo]Page() {
  // 🎛️ Estados locais para UI
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedFilter, setSelectedFilter] = useState('');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<Entity | null>(null);
  const [selected, setSelected] = useState<Entity | null>(null);

  // 🔗 Hook customizado para estado e API
  const {
    items,
    loading,
    error,
    total,
    totalPages,
    create,
    update,
    delete: remove,
    refetch
  } = use[Modulo]({
    search: searchTerm || undefined,
    filter: selectedFilter || undefined,
  });

  // 🎯 Handlers padronizados
  const handleCreate = () => { /* ... */ };
  const handleEdit = (item: Entity) => { /* ... */ };
  const handleDelete = async (id: string) => { /* ... */ };
  const handleView = (item: Entity) => { /* ... */ };

  // 🎨 Estados de loading e error
  if (loading === 'loading') return <LoadingSkeleton />;
  if (error) return <ErrorState error={error} onRetry={refetch} />;

  // 🖼️ Renderização da interface
  return (
    <DashboardLayout>
      <div className="p-6 space-y-6">
        {/* Header, Stats, Filters, Content, Modals */}
      </div>
    </DashboardLayout>
  );
}
```

### 2. Custom Hook Pattern

```typescript
// Padrão para hooks customizados
export function use[Modulo](params: Use[Modulo]Params = {}) {
  // Estados padronizados
  const [items, setItems] = useState<Entity[]>([]);
  const [loading, setLoading] = useState<LoadingState>('idle');
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);

  // Função de fetch
  const fetchItems = async () => {
    setLoading('loading');
    setError(null);
    try {
      const response = await [modulo]Service.getAll(params);
      setItems(response.data);
      setTotal(response.total);
      setTotalPages(response.totalPages);
      setLoading('success');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao carregar');
      setLoading('error');
    }
  };

  // Operações CRUD
  const create = async (data: CreateEntity) => { /* ... */ };
  const update = async (id: string, data: UpdateEntity) => { /* ... */ };
  const remove = async (id: string) => { /* ... */ };

  // Effect para sincronização
  useEffect(() => {
    fetchItems();
  }, [params.search, params.filter]);

  // Retorno padronizado
  return {
    items,
    loading,
    error,
    total,
    totalPages,
    refetch: fetchItems,
    create,
    update,
    delete: remove,
  };
}
```

### 3. Service Layer Pattern

```typescript
// Padrão para services
export class [Modulo]Service {
  private endpoint = '/[modulo]';

  async getAll(params?: GetAllParams): Promise<PaginatedResponse<Entity>> {
    const searchParams = new URLSearchParams();
    
    // Construir parâmetros de busca
    if (params?.search) searchParams.set('search', params.search);
    if (params?.filter) searchParams.set('filter', params.filter);
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.limit) searchParams.set('limit', params.limit.toString());

    const query = searchParams.toString();
    const endpoint = query ? `${this.endpoint}?${query}` : this.endpoint;
    
    const response = await apiClient.get<PaginatedResponse<Entity>>(endpoint);
    return response.data;
  }

  async getById(id: string): Promise<Entity> {
    const response = await apiClient.get<Entity>(`${this.endpoint}/${id}`);
    return response.data;
  }

  async create(data: CreateEntity): Promise<Entity> {
    const response = await apiClient.post<Entity>(this.endpoint, data);
    return response.data;
  }

  async update(id: string, data: UpdateEntity): Promise<Entity> {
    const response = await apiClient.put<Entity>(`${this.endpoint}/${id}`, data);
    return response.data;
  }

  async delete(id: string): Promise<void> {
    await apiClient.delete(`${this.endpoint}/${id}`);
  }
}

export const [modulo]Service = new [Modulo]Service();
```

### 4. Component Patterns

#### UI Components (Reusable)
```typescript
// /components/ui/Button.tsx
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  children: React.ReactNode;
}

export function Button({ variant = 'primary', size = 'md', loading, children, ...props }: ButtonProps) {
  // Implementação padronizada
}
```

#### Shared Components (Cross-module)
```typescript
// /components/shared/DataTable.tsx
interface DataTableProps<T> {
  data: T[];
  columns: TableColumn<T>[];
  loading?: boolean;
  sortState?: SortState;
  onSort?: (field: string) => void;
  emptyMessage?: string;
}

export function DataTable<T>({ data, columns, loading, ...props }: DataTableProps<T>) {
  // Implementação genérica
}
```

#### Domain Components (Module-specific)
```typescript
// /components/domain/produtos/ProdutoCard.tsx
interface ProdutoCardProps {
  produto: Produto;
  onView: (produto: Produto) => void;
  onEdit: (produto: Produto) => void;
  onDelete: (id: string) => void;
}

export function ProdutoCard({ produto, onView, onEdit, onDelete }: ProdutoCardProps) {
  // Implementação específica
}
```

## Data Models

### Entity Interfaces Pattern

```typescript
// /types/entities.ts
export interface BaseEntity {
  id: string;
  dataCadastro: Date;
  ativo: boolean;
}

export interface Produto extends BaseEntity {
  nome: string;
  descricao?: string;
  preco: number;
  categoria: string;
  subcategoria?: string;
  marca: string;
  estoque: number;
  estoqueMinimo: number;
}

export interface Cliente extends BaseEntity {
  nome: string;
  email: string;
  telefone: string;
  endereco?: string;
}
```

### API Response Types

```typescript
// /types/api.ts
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ApiError {
  message: string;
  code?: string;
  details?: any;
}
```

## Error Handling

### Standardized Error Handling

```typescript
// Padrão para tratamento de erros
try {
  const result = await operation();
  return result;
} catch (error) {
  // Log do erro
  console.error(`Erro em ${operationName}:`, error);
  
  // Tratamento baseado no tipo
  if (error instanceof ApiError) {
    throw new Error(error.message);
  }
  
  if (error instanceof NetworkError) {
    throw new Error('Erro de conexão. Verifique sua internet.');
  }
  
  // Erro genérico
  throw new Error('Ocorreu um erro inesperado');
}
```

### Error Boundaries

```typescript
// Padrão para error boundaries
export function ErrorBoundary({ children, fallback }: ErrorBoundaryProps) {
  return (
    <ErrorBoundaryComponent
      fallback={({ error, retry }) => (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800">{error.message}</p>
          <button onClick={retry} className="mt-2 btn-primary">
            Tentar novamente
          </button>
        </div>
      )}
    >
      {children}
    </ErrorBoundaryComponent>
  );
}
```

## Testing Strategy

### Unit Testing Pattern

```typescript
// Padrão para testes de hooks
describe('useProdutos', () => {
  it('should load produtos on mount', async () => {
    const { result } = renderHook(() => useProdutos());
    
    expect(result.current.loading).toBe('loading');
    
    await waitFor(() => {
      expect(result.current.loading).toBe('success');
      expect(result.current.produtos).toHaveLength(3);
    });
  });
});

// Padrão para testes de componentes
describe('ProdutoCard', () => {
  it('should render produto information', () => {
    render(<ProdutoCard produto={mockProduto} onView={jest.fn()} />);
    
    expect(screen.getByText(mockProduto.nome)).toBeInTheDocument();
    expect(screen.getByText(formatCurrency(mockProduto.preco))).toBeInTheDocument();
  });
});
```

## Performance Considerations

### Optimization Patterns

```typescript
// Debounce para busca
const debouncedSearch = debounce((value: string) => {
  setSearchTerm(value);
}, 300);

// Memoização de componentes pesados
const MemoizedComponent = React.memo(ExpensiveComponent);

// Lazy loading de componentes
const LazyComponent = lazy(() => import('./HeavyComponent'));

// Virtualização para listas grandes
import { FixedSizeList as List } from 'react-window';
```

## Security Best Practices

### Input Validation

```typescript
// Validação no frontend
const validateInput = (data: FormData) => {
  const errors = validateForm(data, {
    nome: [required, minLength(2), maxLength(100)],
    email: [required, email],
    preco: [required, numeric, min(0)],
  });
  
  return errors;
};

// Sanitização de dados
const sanitizeInput = (input: string) => {
  return input.trim().replace(/[<>]/g, '');
};
```

### API Security

```typescript
// Headers de segurança
const apiClient = new ApiClient({
  headers: {
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
  },
});

// Validação de responses
const validateResponse = (response: any) => {
  if (!response || typeof response !== 'object') {
    throw new Error('Invalid response format');
  }
  
  return response;
};
```