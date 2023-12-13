import React, { useCallback, useState } from "react";

import Form from "react-bootstrap/Form";

import { UsrCanetDeviceConverterConfig } from "../models/usr-canet-converter-config";
import { Button, Col, Row } from "react-bootstrap";
import { emptyValueOnUndefined } from "../utils/fields";

interface DeviceConfigFormProps {
  value: UsrCanetDeviceConverterConfig;
  onDelete: (_: UsrCanetDeviceConverterConfig) => void;
  onChange: (_: UsrCanetDeviceConverterConfig) => void;
}

const DeviceConfigForm: React.FC<DeviceConfigFormProps> = ({
  value,
  onDelete,
  onChange,
}) => {
  const [editingItem, setEditingItem] = useState(value);

  const onHostChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newItem = { ...editingItem };
      newItem.host = e.target.value;
      setEditingItem(newItem);
      onChange(newItem);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    []
  );

  const onPortChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newItem = { ...editingItem };
      const newPort = parseInt(e.target.value);
      newItem.port = !isNaN(newPort) ? newPort : undefined;
      setEditingItem(newItem);
      onChange(newItem);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [editingItem]
  );

  const onTargetChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newItem = { ...editingItem };
      newItem.target = e.target.value;
      setEditingItem(newItem);
      onChange(newItem);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [editingItem]
  );

  const onDeleteClick = useCallback(() => {
    onDelete(editingItem);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editingItem]);

  return (
    <div className="section border rounded mb-2 p-2">
      <Form.Group>
        <Row className="mb-1">
          <Col md={2}>
            <Form.Label>
              <span className="text-danger">*</span> Host:
            </Form.Label>
          </Col>
          <Col md={2}>
            <Form.Control
              onChange={onHostChange}
              aria-placeholder="Host"
              value={emptyValueOnUndefined(editingItem.host)}
            ></Form.Control>
          </Col>
        </Row>
        <Row className="mb-1">
          <Col md={2}>
            <Form.Label>
              <span className="text-danger">*</span> Port:
            </Form.Label>
          </Col>
          <Col md={2}>
            <Form.Control
              onChange={onPortChange}
              aria-placeholder="Port"
              value={emptyValueOnUndefined(editingItem.port)}
            ></Form.Control>
          </Col>
        </Row>
        <Row className="mb-1">
          <Col md={2}>
            <Form.Label>
              <span className="text-danger">*</span> Target:
            </Form.Label>
          </Col>
          <Col md={2}>
            <Form.Control
              onChange={onTargetChange}
              aria-placeholder="Target"
              value={emptyValueOnUndefined(editingItem.target)}
            ></Form.Control>
          </Col>
        </Row>
        <Row className="mt-2">
          <Col>
            <Button
              className="btn btn-danger float-end"
              onClick={onDeleteClick}
            >
              <i className="fa fa-trash"></i>
            </Button>
          </Col>
        </Row>
      </Form.Group>
    </div>
  );
};

const propsAreEqual = (
  prevItemProps: DeviceConfigFormProps,
  nextItemProps: DeviceConfigFormProps
): boolean => {
  return (
    prevItemProps.value.host === nextItemProps.value.host &&
    prevItemProps.value.port === nextItemProps.value.port &&
    prevItemProps.value.target === nextItemProps.value.target
  );
};

export default React.memo(DeviceConfigForm, propsAreEqual);
