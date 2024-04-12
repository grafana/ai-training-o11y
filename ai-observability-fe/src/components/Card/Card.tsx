

const getStyles = (theme: GrafanaTheme2) => ({
  selectElement: css`
    position: relative;
    width: 225px;
    min-width: 230px;
    flex-grow: 1;
    padding: 24px 30px 24px 24px;
    background: ${theme.colors.background.primary};
    border: 1px solid ${theme.colors.border.medium};
    margin-right: 16px;
    cursor: pointer;
    box-shadow: 0 2px 9px ${theme.isLight ? `rgba(208, 208, 208, 0.6)` : `rgba(0, 0, 0, 0.6)`};
    transition: background 0.2s ease-in-out;

    &:hover {
      background: ${theme.colors.background.secondary};
      box-shadow: inset 0 0 1px ${colors.black};
    }
  `,
  mediumSize: css`
    width: 230px;
    height: 104px;
  `,
  selected: css`
    padding: 23px 30px 23px 23px;
    border: 2px solid ${theme.isLight ? `#5794F2` : `#2f5ca1`};
    box-shadow: 0 2px 9px ${theme.isLight ? colors.blue09 : colors.blue10};
    background: ${theme.colors.background.secondary};

    &:hover {
      box-shadow: 0 2px 9px ${theme.isLight ? colors.blue09 : colors.blue10};
    }
  `,
  icon: css`
    width: 24px;
    height: auto;
    margin-right: 12px;
  `,
  titleWrapper: css`
    display: flex;
    align-items: center;
    margin-bottom: 15px;
  `,
  title: css`
    font-size: ${theme.typography.body.fontSize};
    font-weight: ${theme.typography.fontWeightRegular};
    color: ${theme.isLight ? colors.blue04 : theme.colors.text.maxContrast};
    margin-bottom: 0;
    line-height: 20px;
  `,
  description: css`
    font-size: ${theme.typography.bodySmall.fontSize};
    margin-bottom: 0;
    line-height: 16px;
    color: ${theme.colors.text.secondary};
  `,
  checkbox: css`
    display: block;
    width: 24px;
    min-width: 24px;
    height: 24px;
    margin-left: 15px;
    position: absolute;
    top: 5px;
    right: 5px;
  `,
  checked: css`
    color: ${colors.black};
    position: relative;
    color: ${colors.blue03};
  `,
});

import React, { } from 'react';
import { css, cx } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Icon, useStyles2 } from '@grafana/ui';
import { colors } from 'utils/consts';

export interface CardElementProps<T = any> {
  isSelected?: boolean;
  onClick?: (value?: T) => void;
  size?: 'md' | 'base';
  description: string;
  img?: string;
  title: string;
  component?: React.ReactNode; // New prop for passing the component directly
}

export const CardElement = ({ isSelected, onClick, size, title, img, description, component }: CardElementProps) => {
  const styles = useStyles2(getStyles);

  const handleClick = () => {
    if (onClick) {
      onClick(title); // Pass the title (or any identifier) of the clicked card
    }
  };

  return (
    <div
      onClick={handleClick}
      tabIndex={0}
      role="button"
      className={cx(styles.selectElement, isSelected && styles.selected, size === 'md' && styles.mediumSize)}
    >
      <div className={styles.titleWrapper}>
        <p className={styles.title}>{title}</p>
      </div>
      {component} {/* Render the passed component */}
      {!component && description && <p className={styles.description}>{description}</p>}
      <span className={styles.checkbox}>
        {isSelected && <Icon className={styles.checked} name="check" size="xl" />}
      </span>
    </div>
  );
};
