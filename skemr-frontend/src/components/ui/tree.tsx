// biome-ignore-all lint/suspicious/noExplicitAny: known

import type { ItemInstance } from "@headless-tree/core";
import { CaretDownIcon } from "@phosphor-icons/react";
import * as React from "react";

import { cn } from "@/lib/utils";

interface TreeContextValue<T = any> {
  indent: number;
  currentItem?: ItemInstance<T>;
  tree?: any;
}

const TreeContext = React.createContext<TreeContextValue>({
  currentItem: undefined,
  indent: 20,
  tree: undefined,
});

function useTreeContext<T = any>() {
  return React.useContext(TreeContext) as TreeContextValue<T>;
}

interface TreeProps extends React.HTMLAttributes<HTMLDivElement> {
  indent?: number;
  tree?: any;
}

function Tree({ indent = 20, tree, className, ...props }: TreeProps) {
  const containerProps =
    tree && typeof tree.getContainerProps === "function"
      ? tree.getContainerProps()
      : {};
  const mergedProps = { ...props, ...containerProps };

  const { style: propStyle, ...otherProps } = mergedProps;

  const mergedStyle = {
    ...propStyle,
    "--tree-indent": `${indent}px`,
  } as React.CSSProperties;

  return (
    <TreeContext.Provider value={{ indent, tree }}>
      <div
        className={cn("flex flex-col", className)}
        data-slot="tree"
        style={mergedStyle}
        {...otherProps}
      />
    </TreeContext.Provider>
  );
}

interface TreeItemProps<
  T = any,
> extends React.HTMLAttributes<HTMLButtonElement> {
  item: ItemInstance<T>;
  indent?: number;
}

function TreeItem<T = any>({
  item,
  className,
  children,
  ...props
}: Omit<TreeItemProps<T>, "indent">) {
  const { indent } = useTreeContext<T>();

  const itemProps = typeof item.getProps === "function" ? item.getProps() : {};

  // Merge event handlers so both the user's handlers and tree's internal
  // handlers (selection, keyboard nav) both fire.
  const {
    onClick: propsOnClick,
    onKeyDown: propsOnKeyDown,
    style: propStyle,
    ...restProps
  } = props;
  const {
    onClick: itemOnClick,
    onKeyDown: itemOnKeyDown,
    style: itemStyle,
    ...restItemProps
  } = itemProps;

  const mergedProps = {
    ...restProps,
    ...restItemProps,
    onClick: (e: React.MouseEvent<HTMLButtonElement>) => {
      propsOnClick?.(e);
      itemOnClick?.(e);
    },
    onKeyDown: (e: React.KeyboardEvent<HTMLButtonElement>) => {
      propsOnKeyDown?.(e);
      itemOnKeyDown?.(e);
    },
  };

  const mergedStyle = {
    ...propStyle,
    ...itemStyle,
    "--tree-padding": `${item.getItemMeta().level * indent}px`,
  } as React.CSSProperties;

  const otherProps = mergedProps;

  return (
    <TreeContext.Provider value={{ currentItem: item, indent }}>
      <button
        aria-expanded={item.isExpanded()}
        className={cn(
          "z-10 select-none ps-(--tree-padding) not-last:pb-0.5 outline-hidden focus:z-20 data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
          className,
        )}
        data-focus={
          typeof item.isFocused === "function"
            ? item.isFocused() || false
            : undefined
        }
        data-folder={
          typeof item.isFolder === "function"
            ? item.isFolder() || false
            : undefined
        }
        data-search-match={
          typeof item.isMatchingSearch === "function"
            ? item.isMatchingSearch() || false
            : undefined
        }
        data-selected={
          typeof item.isSelected === "function"
            ? item.isSelected() || false
            : undefined
        }
        data-slot="tree-item"
        style={mergedStyle}
        {...otherProps}
      >
        {children}
      </button>
    </TreeContext.Provider>
  );
}

interface TreeItemLabelProps<
  T = any,
> extends React.HTMLAttributes<HTMLSpanElement> {
  item?: ItemInstance<T>;
}

function TreeItemLabel<T = any>({
  item: propItem,
  children,
  className,
  ...props
}: TreeItemLabelProps<T>) {
  const { currentItem } = useTreeContext<T>();
  const item = propItem || currentItem;

  if (!item) {
    console.warn("TreeItemLabel: No item provided via props or context");
    return null;
  }

  return (
    <span
      className={cn(
        "flex items-center gap-1 rounded-sm bg-background in-data-[drag-target=true]:bg-accent in-data-[search-match=true]:bg-blue-400/20! in-data-[selected=true]:bg-accent px-1.5 py-0.5 not-in-data-[folder=true]:ps-7 in-data-[selected=true]:text-accent-foreground text-xs in-focus-visible:ring-[3px] in-focus-visible:ring-ring/50 transition-colors hover:bg-accent [&_svg]:pointer-events-none [&_svg]:shrink-0",
        className,
      )}
      data-slot="tree-item-label"
      {...props}
    >
      {item.isFolder() && (
        <CaretDownIcon className="in-aria-[expanded=false]:-rotate-90 size-4 text-muted-foreground transition-transform" />
      )}
      {children ??
        (typeof item.getItemName === "function" ? item.getItemName() : null)}
    </span>
  );
}

export { Tree, TreeItem, TreeItemLabel };
